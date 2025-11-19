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

---

# 心得
在閱讀一本高併發系統設計相關的書籍時，其中提到「點讚服務（Like Service）」的實作方式。基於好奇心，我決定自行開發一個點讚服務系統，並嘗試探索不同架構對高併發的處理能力


起初我以 MySQL 作為資料儲存層，並採用悲觀鎖來處理併發寫入問題。然而在高頻率的 UPDATE操作下，仍然對 MySQL 造成了明顯的性能負擔，導致吞吐量受限


為改善這個瓶頸，我改以 Redis 作為快取與計數層，雖然能有效減少 MySQL 的寫入壓力，但直接在 Redis 中頻繁修改計數仍然會造成資源浪費，且無法完全避免大量併發帶來的成本


因此，我進一步將架構調整為以 Kafka 作為中間層（Message Queue）。

 透過 Kafka 的高吞吐量特性，以及 Consumer Group 的併行消費能力，系統可以：

    先累積一段時間或一定數量的點讚事件
    再批次寫入 MySQL 或其他儲存層
    避免高併發下直接打爆資料庫
    大幅降低重複寫入的成本


在效能監控方面，我整合 Prometheus + Grafana，用於可視化 Kafka、Redis、API 層的性能指標，包含 QPS、延遲、錯誤率等


為提升 API 層的穩定性，我實作了 指數回退（Exponential Backoff）重試機制，以提高在短暫錯誤情境下的成功率，並減少對後端的不必要壓力
