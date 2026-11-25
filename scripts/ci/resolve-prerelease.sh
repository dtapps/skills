#!/usr/bin/env bash
#
# 版本号归一化 + 预发布判定。
#
# 事实来源：GitHub Actions 的 release.yml / build-image.yml 与 CNB 的
# .cnb/workflows/release.yml / build-image.yml 都 source 这一份实现，避免四处各写一套
# grep 正则、判断口径慢慢跑偏：
#
#   v1.2.3 / 1.2.3        -> TAG=v1.2.3        PRERELEASE=false   正式版
#   v1.2.3-rc.1           -> TAG=v1.2.3-rc.1   PRERELEASE=true    预发布
#   v1.2.3-beta.1         -> TAG=v1.2.3-beta.1 PRERELEASE=true    预发布
#   v1.2.3-alpha.1        -> TAG=v1.2.3-alpha.1 PRERELEASE=true   预发布
#   v1.2.3-nightly.20260101 -> TAG=v1.2.3-nightly.20260101  PRERELEASE=true  预发布（nightly 渠道）
#   nightly               -> 自动展开为 v<最近正式版+1补丁>-nightly.<YYYYMMDD>
#
# 只用版本号后缀判定，不接受「是不是预发布」这种布尔输入：布尔输入会和版本号打架。
#
# nightly 的两种写法，最终都落到带日期的合法 semver 上：
#   - 输入写 `nightly`：取最近正式版 +1 补丁做基版本（1.2.3 -> 1.2.4），
#     日期取 UTC 当天，生成 v1.2.4-nightly.20260101。
#     每天一个 tag；发布流水线发布后只保留最近 3 个 nightly Release。
#     基版本必须比已发布的正式版高，否则 semver 会把 1.2.3 判成比
#     1.2.3-nightly.X 新，更新器会误提示「升级到旧正式版」。
#     没有正式版 tag 时（备用远端可能只有 rc tag）退而取最大的预发布基版本，
#     且不再 +1：rc/nightly 的 X.Y.Z 已经是待发布版本，再 +1 会跳过补丁号。
#   - 输入直接写日期版本（v1.2.4-nightly.20260101）：按普通预发布处理。
#
# 用法（被 source）：
#   . scripts/ci/resolve-prerelease.sh
#   resolve_version "${VERSION}"
#   printf '%s %s\n' "${TAG}" "${PRERELEASE}"
#
# 自检可用环境变量固定 nightly 的展开结果（CI 不需要设置）：
#   NIGHTLY_BASE_VERSION  基版本（X.Y.Z），替代自动探测的最近正式版
#   NIGHTLY_DATE          日期（YYYYMMDD），替代 UTC 当天
#
# 直接执行可以跑内置用例自检（改了正则以后先跑它）：
#   bash scripts/ci/resolve-prerelease.sh --self-test
#
# ⚠️ 本文件必须保持 POSIX sh 语法：CNB 的 stage 默认用 `sh`（可能是 dash/ash）执行，
# GitHub Actions 用 bash。不要引入 [[ ]]、数组、BASH_SOURCE、$'...' 这类 bash 专有写法。

# resolve_version <版本号>
#   归一化版本号（补 v 前缀）并判定是否预发布；
#   结果写入全局变量 TAG / PRERELEASE / VERSION_TYPE / CHANNEL / VERSION_BARE / PRE_NOTE
#   并导出给后续步骤：
#     TAG          git tag / Release 用的版本号，如 v1.2.3-rc.1 或 v1.2.4-nightly.20260101
#     PRERELEASE   true / false
#     VERSION_TYPE 正式版 / 预发布
#     CHANNEL      stable / rc / beta / alpha / nightly，供 changelog 选对比基准
#     VERSION_BARE 去掉 v 前缀，供展示与通知使用
#   版本号为空时返回 1，由调用方决定是 `set -e` 直接失败还是自行处理。
resolve_version() {
  local input="${1:-}"

  # 去掉所有空白字符：workflow_dispatch 输入框粘贴时经常带上首尾空格/换行
  input="$(printf '%s' "${input}" | tr -d '[:space:]')"

  if [ -z "${input}" ]; then
    printf '错误: 版本号不能为空。示例: v1.2.3 / v1.2.3-rc.1 / v1.2.3-beta.1 / v1.2.3-alpha.1 / v1.2.3-nightly.20260101 / nightly\n' >&2
    return 1
  fi

  # 滚动预发布（nightly）：展开成带日期的合法 semver。
  #   基版本 = 最近正式版 +1 补丁（1.2.3 -> 1.2.4），日期 = UTC 当天（YYYYMMDD）。
  #   没有正式版 tag 时，取最大的预发布基版本直接用（不再 +1）。
  # 基版本必须高于已发布正式版：semver 里 1.2.3-nightly.X < 1.2.3，
  # 用原版本号做基版本会让更新器误提示「升级到旧正式版」。
  #
  # 必须在「补 v 前缀」之前判定：nightly 不是版本号，
  # 补前缀后会变成 vnightly，被当成普通预发布 tag 发出去。
  case "${input}" in
    nightly | NIGHTLY)
      local base date_part major rest minor patch next bump
      base="${NIGHTLY_BASE_VERSION:-}"
      bump='yes'
      if [ -z "${base}" ]; then
        # 只认严格的三段式版本 tag（vX.Y.Z）。
        # 不能只排除「-rc / -beta / -alpha / -nightly」这类连字符后缀：
        # 历史上存在 v1.0.44.rc1 这种「点分隔」的预发布 tag，它会混进来被当成
        # 最近正式版，导致基版本不是 X.Y.Z、nightly 推导直接失败。
        base="$(git tag -l 'v*' 2>/dev/null \
          | grep -E '^[vV]?[0-9]+\.[0-9]+\.[0-9]+$' \
          | sed 's/^[vV]//' | sort -V | tail -n 1 || true)"
      fi
      if [ -z "${base}" ]; then
        # 备用远端（如 GitHub）可能只同步过 rc tag，一个正式版 tag 都没有。
        # 这时取预发布 tag 的 X.Y.Z 直接用：它们已经是「待发布版本」，
        # 再 +1 会跳过一个补丁号，并和主平台的 nightly 版本号对不上。
        base="$(git tag -l 'v*' 2>/dev/null \
          | grep -E '^[vV]?[0-9]+\.[0-9]+\.[0-9]+-(rc|beta|alpha|nightly)([.-][0-9A-Za-z.-]*)?$' \
          | sed -E 's/^[vV]([0-9]+\.[0-9]+\.[0-9]+)-(rc|beta|alpha|nightly).*$/\1/' \
          | sort -V | tail -n 1 || true)"
        bump='no'
      fi
      base="${base:-0.0.0}"
      if ! printf '%s' "${base}" | grep -Eq '^[0-9]+\.[0-9]+\.[0-9]+$'; then
        printf '错误: nightly 基版本不是 X.Y.Z: %s\n' "${base}" >&2
        return 1
      fi
      major="${base%%.*}"
      rest="${base#*.}"
      minor="${rest%%.*}"
      patch="${rest#*.}"
      if [ "${bump}" = 'yes' ]; then
        next="${major}.${minor}.$((patch + 1))"
      else
        next="${base}"
      fi

      date_part="${NIGHTLY_DATE:-$(date -u +%Y%m%d)}"
      if ! printf '%s' "${date_part}" | grep -Eq '^[0-9]{8}$'; then
        printf '错误: nightly 日期不是 YYYYMMDD: %s\n' "${date_part}" >&2
        return 1
      fi

      TAG="v${next}-nightly.${date_part}"
      PRERELEASE="true"
      VERSION_TYPE="预发布"
      CHANNEL="nightly"
      VERSION_BARE="${next}-nightly.${date_part}"
      PRE_NOTE="版本号: ${TAG}（${VERSION_TYPE}，prerelease=${PRERELEASE}，channel=nightly，基版本 ${next}）"
      export TAG PRERELEASE VERSION_TYPE CHANNEL VERSION_BARE PRE_NOTE
      return 0
      ;;
  esac

  # 字符白名单：版本号最终会变成 git tag、Docker tag 和 JSON 字段值。
  # 含引号/反斜杠/换行的输入既不是合法版本号，拼进 JSON 时还会破坏结构（甚至注入字段），
  # 所以在这里一次性挡掉。允许的字符与 Docker tag 的字符集一致。
  if printf '%s' "${input}" | grep -qvE '^[0-9A-Za-z._+-]+$'; then
    printf '错误: 版本号含非法字符（只允许字母、数字、. _ + -）: %s\n' "${input}" >&2
    return 1
  fi

  # 统一补 v 前缀：输入 1.2.3 与 v1.2.3 等价；大写 V 前缀也算已经有了，
  # 否则会拼出 vV1.2.3 这种既不像 tag 又不像版本号的东西
  case "${input}" in
    v* | V*) TAG="${input}" ;;
    *) TAG="v${input}" ;;
  esac

  # 后缀判定：-rc / -beta / -alpha / -nightly（大小写不敏感，rc1 / rc.1 / 裸 rc 都算预发布）
  # nightly 后缀带日期（如 -nightly.20260101），和另外三个渠道一样只看版本号末尾。
  # 匹配要求后缀出现在版本号末尾，避免 v1.2.3-source、v1.2.3-hotfix 这类被误判
  if printf '%s' "${TAG}" | grep -qiE -- '-(rc|beta|alpha|nightly)([.-]?[0-9A-Za-z.-]*)$'; then
    PRERELEASE="true"
    VERSION_TYPE="预发布"
  else
    PRERELEASE="false"
    VERSION_TYPE="正式版"
  fi

  # 渠道判定：后缀决定渠道，changelog 按渠道选对比基准（rc 对 rc、正式版对正式版）。
  # 只认 -rc / -beta / -alpha / -nightly 四种后缀，且必须在末尾。
  case "${TAG}" in
    *-rc*) CHANNEL="rc" ;;
    *-beta*) CHANNEL="beta" ;;
    *-alpha*) CHANNEL="alpha" ;;
    *-nightly*) CHANNEL="nightly" ;;
    *) CHANNEL="stable" ;;
  esac

  # 展示用版本号不带 v 前缀（通知 / Release 正文里用），大写 V 前缀同样要剥掉，
  # 否则会拼出 V1.2.3 这种两不像的值。
  case "${TAG}" in
    [vV]*) VERSION_BARE="${TAG#?}" ;;
    *) VERSION_BARE="${TAG}" ;;
  esac

  PRE_NOTE="版本号: ${TAG}（${VERSION_TYPE}，prerelease=${PRERELEASE}，channel=${CHANNEL}）"
  export TAG PRERELEASE VERSION_TYPE CHANNEL VERSION_BARE PRE_NOTE

  return 0
}

# 下面是「直接执行」才跑的自检，被 source 时只导出函数就收工。
#
# 判定方式不能用 `[ "${BASH_SOURCE[0]}" = "${0}" ]`：数组下标是 bash 专有语法，CNB 的 stage
# 默认用 sh（dash）执行，`${BASH_SOURCE[0]:-}` 会直接 "Bad substitution" 让整个 stage 挂掉。
# 这里改用 $0 的 basename：直接执行时 $0 是本文件名，被 source 时 $0 是调用方
# （bash / sh / /bin/dash 或调用脚本），在 POSIX 范围内可靠。
# 另外放行显式的 --self-test：source 的调用方不会传这个参数，不会误触发。
if [ "${0##*/}" != 'resolve-prerelease.sh' ] && [ "${1:-}" != '--self-test' ]; then
  # `return` 在被 source 时生效；万一本文件被 `bash < file` 之类方式执行（$0 不是文件名、
  # return 会失败），退回 exit 0。shellcheck 认为 return 不会失败，故显式关掉该提示。
  # shellcheck disable=SC2317
  return 0 2>/dev/null || exit 0
fi

if [ "${1:-}" != "--self-test" ]; then
  printf '用法: %s --self-test\n通常不需要单独执行，CI 里是 source 使用：%s\n' "${0}" '. scripts/ci/resolve-prerelease.sh' >&2
  exit 2
fi

failures=0

check_version() {
  local input="$1" want_tag="$2" want_pre="$3"

  if ! resolve_version "${input}" 2>/dev/null; then
    printf '  ✗ %-20s 解析失败（期望 %s）\n' "${input:-<空>}" "${want_tag}"
    failures=$((failures + 1))
    return 0
  fi

  if [ "${TAG}" = "${want_tag}" ] && [ "${PRERELEASE}" = "${want_pre}" ]; then
    printf '  ✓ %-20s -> %-20s prerelease=%s\n' "${input}" "${TAG}" "${PRERELEASE}"
  else
    printf '  ✗ %-20s -> %-20s prerelease=%s（期望 %s / %s）\n' \
      "${input}" "${TAG}" "${PRERELEASE}" "${want_tag}" "${want_pre}"
    failures=$((failures + 1))
  fi
}

check_version 'v1.2.3' 'v1.2.3' 'false'
check_version '1.2.3' 'v1.2.3' 'false'
check_version 'v1.2.3-rc.1' 'v1.2.3-rc.1' 'true'
check_version 'v1.2.3-beta.1' 'v1.2.3-beta.1' 'true'
check_version 'v1.2.3-alpha.1' 'v1.2.3-alpha.1' 'true'
check_version 'v1.2.3-nightly.20260101' 'v1.2.3-nightly.20260101' 'true'
check_version '1.2.3-nightly.20260101' 'v1.2.3-nightly.20260101' 'true'
check_version 'V1.2.3-NIGHTLY.20260101' 'V1.2.3-NIGHTLY.20260101' 'true'
check_version 'v1.2.3-nightly' 'v1.2.3-nightly' 'true'
check_version 'v1.2.3-rc' 'v1.2.3-rc' 'true'
check_version 'V1.2.3-RC.1' 'V1.2.3-RC.1' 'true'
check_version 'v1.2.3-alpha.1.2' 'v1.2.3-alpha.1.2' 'true'
check_version ' v1.2.3 ' 'v1.2.3' 'false'
# nightly：输入 nightly 展开为「最近正式版 +1 补丁 -nightly.UTC日期」。
# 用环境变量固定基版本与日期，避免自检依赖仓库 tag 和当天日期。
NIGHTLY_BASE_VERSION='1.2.3'
NIGHTLY_DATE='20260101'
export NIGHTLY_BASE_VERSION NIGHTLY_DATE
check_version 'nightly' 'v1.2.4-nightly.20260101' 'true'
check_version 'NIGHTLY' 'v1.2.4-nightly.20260101' 'true'
# 补丁号进位：1.2.9 -> 1.2.10
NIGHTLY_BASE_VERSION='1.2.9'
check_version 'nightly' 'v1.2.10-nightly.20260101' 'true'
NIGHTLY_BASE_VERSION='1.2.3'
# 不带预发布后缀（rc / beta / alpha / nightly）的一律是正式版
check_version 'v1.2.3-hotfix' 'v1.2.3-hotfix' 'false'
check_version 'v2.0.0-core' 'v2.0.0-core' 'false'
# 没有 v 前缀的裸 nightly 不能被当成普通版本号（补前缀会变成 vnightly）
check_version 'v1.2.3-night' 'v1.2.3-night' 'false'

# VERSION_BARE 用于展示与通知，必须去掉 v 前缀
check_bare() {
  local input="$1" want="$2"
  resolve_version "${input}" >/dev/null 2>&1 || {
    printf '  ✗ %-20s 解析失败\n' "${input}"
    failures=$((failures + 1))
    return 0
  }
  if [ "${VERSION_BARE}" = "${want}" ]; then
    printf '  ✓ %-20s -> VERSION_BARE=%s\n' "${input}" "${VERSION_BARE}"
  else
    printf '  ✗ %-20s -> VERSION_BARE=%s（期望 %s）\n' "${input}" "${VERSION_BARE}" "${want}"
    failures=$((failures + 1))
  fi
}
check_bare 'v1.2.3' '1.2.3'
check_bare '1.2.3' '1.2.3'
check_bare 'V1.2.3-RC.1' '1.2.3-RC.1'
check_bare 'v1.2.3-nightly.20260101' '1.2.3-nightly.20260101'
check_bare 'nightly' '1.2.4-nightly.20260101'

# 渠道要能和 changelog 的口径对上
check_channel() {
  local input="$1" want="$2"
  resolve_version "${input}" >/dev/null 2>&1 || {
    printf '  ✗ %-20s 解析失败\n' "${input}"
    failures=$((failures + 1))
    return 0
  }
  if [ "${CHANNEL}" = "${want}" ]; then
    printf '  ✓ %-20s -> channel=%s\n' "${input}" "${CHANNEL}"
  else
    printf '  ✗ %-20s -> channel=%s（期望 %s）\n' "${input}" "${CHANNEL}" "${want}"
    failures=$((failures + 1))
  fi
}
check_channel 'v1.2.3' 'stable'
check_channel 'v1.2.3-rc.1' 'rc'
check_channel 'v1.2.3-beta.2' 'beta'
check_channel 'v1.2.3-alpha.3' 'alpha'
check_channel 'v1.2.3-nightly.20260101' 'nightly'
check_channel 'nightly' 'nightly'
unset NIGHTLY_BASE_VERSION NIGHTLY_DATE

# 非法输入必须被拒绝，不能漏到 git tag / Docker tag / JSON 里：
# 空版本号历史上会打出 :v 这种非法标签；引号/分号会破坏 optionsFrom 的 JSON 结构。
# 注意不要把 "v1.2.3 beta" 当非法输入：空白会先被去掉（粘贴常带空格），结果 v1.2.3beta 合法
for bad_input in '' 'v1.2.3"x' 'v1.2.3;rm -rf /' 'v1.2.3@1' 'v1/2/3'; do
  if resolve_version "${bad_input}" 2>/dev/null; then
    printf '  ✗ %-20s 应当报错，实际通过（TAG=%s）\n' "${bad_input:-<空>}" "${TAG}"
    failures=$((failures + 1))
  else
    printf '  ✓ %-20s 被拒绝（返回非 0）\n' "${bad_input:-<空>}"
  fi
done

if [ "${failures}" -ne 0 ]; then
  printf '\n自检失败: %d 项\n' "${failures}" >&2
  exit 1
fi

printf '\n自检通过\n'
