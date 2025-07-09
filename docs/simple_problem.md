# Simple problem
## Chạy chương trình
Cập nhật input trong file test/simple_input.txt
```
KNC ETH
2
KNC USDT 1.1 0.9
ETH USDT 360 355
```

Chạy lệnh
```
go run cmd/simple/main.go
```

## Cách giải quyết
### Mô hình hóa bài toán
Mô hình hóa bài toán theo hướng graph:

Coi mỗi loại `currency` là một `đỉnh` trong graph. Theo đó, mỗi `trading pair` tương ứng với hai `cạnh` :

- Cạnh 1: chiều từ đỉnh 1 `base currency` sang đỉnh 2 `quote currency`. Trọng số tương ứng sẽ là `bid price` và `ask price` 
- Cạnh 2: chiều từ đỉnh 2 `quote currency` sang đỉnh 1 `base currency`. Trọng số tương ứng là `1 / ask price` và `1 / bid price` do đảo chiều

Đối với input của đề bài ta có thể thể hiện dạng graph như sau:
```
KNC ETH
2
KNC USDT 1.1 0.9
ETH USDT 360 355
```
![Simple graph](images/simple_graph.PNG)

### Giải pháp

Yêu cầu đề bài:

1. Best Ask Price (Buying the Base Currency): xác định trading route và giá thấp nhất để mua 1 `base currency` → tìm route từ `quote currency` về `base currency` và nhân trọng số của các cạnh, kết quả sẽ ra được 1 `quote currency` mua được bao nhiêu `base currency` , giá trị này cần phải càng nhỏ càng tốt (vì cần mua với giá thấp) → Do cần tìm route sao cho tích các trọng số (rate) nhỏ nhất, và các trọng số (rate) có thể nằm trong khoảng (0, 1) nên có thể dẫn đến vòng lặp vô hạn (chu trình âm). Vì vậy, không dùng được Dijkstra, cần dùng  **Bellman-Ford**, với trọng số là `log(rate)`, và phép cộng để thay cho phép nhân. 
2. **Best Bid Price (Selling the Base Currency):** xác định trading route và giá cao nhất để bán 1 `base currency` → tìm route từ `base currency` về `quote currency` sao cho tích trọng số của các cạnh là lớn nhất (vì cần bán giá cao nhất) → Tương tự lập luận ở trên nhưng dùng `-log(rate)` để đảo dấu trọng số

Kết luận: bài toán best ask price là bài tìm đường đi ngắn nhất từ base currency về quote currency, bài toán best bid price là bài toán tìm đường đi dài nhất nhưng theo hướng ngược lại.