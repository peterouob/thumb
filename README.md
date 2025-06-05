# 一個計數服務

| name | router|
|------|-------|
|thumb|模擬按讚動作|
|metrics|查看prometheus訊|

# 功能
1. 當使用者請求thumb時將佔存在一個Array中
2. 當用量到一定以後寫入redis並改變offset
3. 假設用量沒到一定的情況但時間到flush time則一樣寫入
4. 假設寫入失敗執行backOff(參考grpc connection backoff)

# 後續優化
- [ ] 限流並且記錄流量
- [ ] 高集中熱點用戶使用RockDB而非redis
- [ ] 將redis寫入mysql保持數據持久化,避免長時間佔用redis