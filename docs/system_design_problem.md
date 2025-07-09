# System design problem
## Step 1: Outline use cases, constraints, and assumptions
### Use cases (Functional Requirements)
Phạm vi hệ thống chỉ xử lý các use cases sau:
- **Service** fetch order book data từ cryptocurrency exchange:
    - Fetch real-time order book data từ exchange.
    - Lưu và cập nhật order book data nếu có thay đổi.
- **User** query best trade route và best trade price:
    - **REST API Endpoint**: input starting `token Y`, target `token X`, trade amount `n` (units of token X). API trả về:
        - The lowest effective ask price when buying n units of Token X và best ask route.
        - The highest effective bid price when selling n units of Token X và best bid route.
    - **Multi-Hop Trading Support**: Consider multiple market pairs if a direct trading pair does not exist.

Out of scope:
- Lưu trữ lâu dài order book phục vụ aggregation về sau
- Trade route tối ưu nằm trên nhiều exchanges

### Constraints and assumptions (Non-Functional Requirements)
Các ràng buộc đưa ra bởi đề bài:
- **Low Latency**: API đưa ra near real-time responses. Order book data cũng được fetch real-time.
- **High Availability**: The system should be resilient to API failures from exchanges.
- **Scalability**: The architecture should handle multiple exchanges and large trading volumes.
- **Fault Tolerance**: Ensure fallback mechanisms in case of API failures or data inconsistencies.
- **Security**: Secure API endpoints and prevent abuse (e.g., rate limiting).

Các con số:

Thông số API của một exchange điển hình (ở đây là Binance):

- Số lượng trading pair (symbol) đang hoạt động trên Binance là 1454 
-> Giả định mỗi exchange có khoảng 1500 symbols hoạt động. Hệ thống liên tục 
fetch và lưu order book của toàn bộ 1500 symbols này.

- API của Binance thực hiện rate limit theo weight, mỗi API đều tính weight. 
Binance giới hạn 6000 weight/IP/phút. Đối với API lấy order book của một 
trading pair, weight tính theo depth:
    ```
    weight = 5, depth 1 - 100
    25, 101 - 500
    50, 501 - 1000
    250, 1001-5000
    ```

- Binance cũng hỗ trợ stream thay đổi của order book qua websocket, mỗi 
kết nối listen tối đa 1024 streams, sau 24h sẽ bị ngắt kết nối. Mỗi IP 
cho phép tạo tối đa 300 connections mỗi 5 phút.

## Step 2: Create a high level design
![High level design](images/high_level_design.drawio.png)

## Step 3: Design core components

### Use case: Service fetch order book data từ cryptocurrency exchange
Giả định thực hiện fetch order book từ 100 exchanges, mỗi exchage theo đề
cập ở trên có khoảng 1500 symbols.

Lưu danh sách exchange và metadata vào bảng `exchanges` (SQL hoặc NoSQL), hoặc file cấu hình, 
do đây là thông tin ít thay đổi.

Với mỗi exchange, Order Book Fetcher Service thực hiện:
- Lấy danh sách symbol đang hoạt động trên exchange
    - Call API của exchange
- Với mỗi symbol, thực hiện lấy order book hiện tại

Open a stream to wss://fstream.binance.com/stream?streams=btcusdt@depth.
Buffer the events you receive from the stream. For same price, latest received update covers the previous one.
Get a depth snapshot from https://fapi.binance.com/fapi/v1/depth?symbol=BTCUSDT&limit=1000 .
Drop any event where u is < lastUpdateId in the snapshot.
The first processed event should have U <= lastUpdateId**AND**u >= lastUpdateId
While listening to the stream, each new event's pu should be equal to the previous event's u, otherwise initialize the process from step 3.
The data in each event is the absolute quantity for a price level.
If the quantity is 0, remove the price level.
Receiving an event that removes a price level that is not in your local order book can happen and is normal.




, nếu fetch 1500 symbol thì mỗi symbol 
nhận được weight = 4, tương đương depth = 80. Đến 1 phút sau mới được fetch 
tiếp, không đảm bảo fetch real-time. Do đó: 
- Chỉ sử dụng API để fetch order book snapshot ban đầu hoặc trường hợp bị 
out-of-sync, sau đó chuyển qua websocket để nhận cập nhật. 
- Nếu cần fetch depth lớn, có thể sử dụng proxy/IP rotation.
- Khi fetch thì nên ưu tiên depth cao cho top trade token, và giảm depth của 
token còn lại.

## Tham khảo
Cách lấy số lượng trading pair (symbol) đang hoạt động trên Binance
```
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)
	exchangeInfo, _ := client.NewExchangeInfoService().SymbolStatus("TRADING").Do(context.Background())
	fmt.Println(len(exchangeInfo.Symbols))
```

[Binance - How to manage a local order book correctly](https://developers.binance.com/docs/derivatives/usds-margined-futures/websocket-market-streams/How-to-manage-a-local-order-book-correctly)

