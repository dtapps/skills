---
name: baidu-baike-data
description: 百度百科组件是一个知识服务工具，用于查询各类名词的权威百科解释。其核心功能是：给定用户提供的特定名词（事物、人物、地点、概念、事件等），从百度百科返回标准化、详细的词条说明。
homepage: https://baike.baidu.com/
metadata: { "openclaw": { "emoji": "📖", "requires": { "bins": ["python3"] } } }
---

# Baidu Baike

Query encyclopedia entries from Baidu Baike.

## Two Usage Scenarios

### Scenario 1: Direct Search
Get default matching entry for a keyword.
```bash
python3 scripts/baidu_baike.py --search_type=lemmaTitle --search_key="keyword"
```

### Scenario 2: Homonym Resolution
When term has multiple entries, list them and select by ID.
```bash
# List entries with same name
python3 scripts/baidu_baike.py --search_type=lemmaList --search_key="keyword" --top_k=5

# Get specific entry by ID
python3 scripts/baidu_baike.py --search_type=lemmaId --search_key="entry_id"
```

## API
- LemmaList: List entries with same title
- LemmaContent: Get entry details by title or ID

## Workflow
1. Extract noun from query
2. For ambiguous terms, call LemmaList first
3. User selects entry from list
4. Call LemmaContent with selected ID
5. Return structured data
