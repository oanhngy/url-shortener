
# MÔ TẢ BÀI TOÁN
Yêu cầu chung:
- Tạo short URL từ long URL
- Truy xuất lại URL gốc khi user truy cập short URL
- Theo dõi lượt click
- Xem URL list đã tạo
- Cần chạy ổn định, mở rộng được và xử lý các vấn đề thường gặp như concurrency, collision, đảm bảo khi lượng URL tăng lớn

Cách hiểu: URL Shortener có 2 phần chính:
1. Generate short URL: Nhận long URL --> tạo short URL --> lưu mapping short-long vào db
2. Redirect: user truy cập short URL --> tìm long URL --> redirect user và tăng click

Scale assumption: Không yêu cầu số lượng URL cần xử lý --> em tự đưa ra các giả định phù hợp với bài test:
- Short code gồm 6 characters
- Sử dụng tập Base62
- Tổng: 62^2 =~ 56 tỷ short URL --> collision tăng nhiều khi số lượng ~238.000 URL --> cần sử lý collision

Chiến lược xử lý collision:
- short code tạo bằng random Base62
- trước khi lưu cần check short_code vừa tạo đã tồn tạo trong db
- db là MongoDB sử dụng unique index trên short code
- Nếu collision: Generate lại - Retry cho đến khi không trùng


# CÁCH CHẠY PROJECT
Môi trường:
- Go >= 1.21
- Docker desktop
- Git

Clone:
git clone https://github.com/oanhngy/url-shortener.git
cd url-shortener

Chạy MongoDB + API server:
docker compose up -d
docker ps
go run ./cmd/api
Server chạy tại: http://localhost:8080

Test bằng curl:
curl -X POST http://localhost:8080/api/v1/links \
  -H "Content-Type: application/json" \
  -d '{"longUrl":"https://example.com/abc"}'

curl http://localhost:8080/api/v1/links

curl http://localhost:8080/api/v1/links/{shortCode}

http://localhost:8080/{shortCode}

Kiểm dữ liệu MongoDB:
docker exec -it urlshortener-mongo mongosh -u root -p rootpassword --authenticationDatabase admin

use urlshortener
db.links.find()


THIẾT KẾ & QUYẾT ĐỊNH KỸ THUẬT
- Database chọn MongoDB vì: 
    - Phù hợp với dữ liệu: chỉ cần lưu mapping short-long code và click count, không cần join, transaction phức tạp
    - Dễ scale: MongoDB phù hợp kiểu nhiều read/write đơn giản, dễ scale khi traffic tăng
    - Xử lý collision/concurrency tốt: tạo unique index trên short_code đảm bảo k có 2 record trùng

- Thiết kế API = REST + net/http:
    - Bài ít endpoint --> REST rõ, dễ test bằng curl
    - net/http + ServeMux tối giản dependency, phù hợp với yêu cầu đề

- Thuật toán generate mã ngắn = Random Base62
    - URL-friendly, không cần thêm decode
    - Độ dài cố định, ngắn

- Xử lý conflict/duplicate: Collision có thể xảy ra khi random trùng mã, 2 requests đồng thời  generate cùng code. Cách xử lý:
    - check tồn tại trước khi insert
    - DB có unique index short_code
    - Bị duplicate --> generate lại + retry

TRADE-OFFS
- chọn Random Base 62 thay vì count-based vì không cần global counter, tránh được bottleneck khi scale; phù hợp với hệ thống phân tán; Đổi lại có khả năng bị collision (đã xử lý) và phụ thuộc vào unique index
- Chọn MongoDB (nonSQL) thay vì SQL lý do: mô hình key-value đơn giản, dễ xử lý unique index + retry cho collision; đổi lại MongoDB không có acid properties như SQL, không phù hợp nếu cần query phức tạp sau này
- net/http + ServeMux thay vì Gin/Fiber: dễ đọc và debug, code minh bạch; đổi lại không có sẵn middleware, binding, validation

CHALLENGES
- Gặp vấn đề gì?
    - Công nghệ mới: mất thời gian làm quen syntax, tolling và cách wiring --> test in-memory trước, thay bằng MongoDB, test từng phần
    - Vấn đề collision/concurrency: random không là chưa đủ --> phải đảm bảo unique constraint (unique index MongoDB), thêm Exists()
    - Wiring toàn bộ backend: Cần hiểu rõ các luồng hoạt động
- Học được gì?
    - Chia nhỏ từng phần để giải quyết
    - Cách tổ chức backend có tổ chức
    - Vai trò của từng layer


LIMITATIONS & IMPROVEMENTS
- Code hiện tại còn thiếu gì?
    - Phần testing ít: chỉ mới unit cho core service, chưa test với MongoDB và HTTP handlers
    - Chưa validate đầy đủ input URL
    - Chưa xử lý các edge cases như URL quá dài...

- Nếu có thêm thời gian sẽ làm gì?
    - Thực hiện testing kỹ hơn (unit test cho MongoDB và HTTP HAndlers)
    - Cải thiện validation (giới hạn độ dài, trả error message rõ ràng)
    - Thêm giao diện web đơn giản (để chạy demo trực quan hơn)

- Production-ready cần thêm gì?
    - Rate limiting chống spam tạo link
    - Monitoring & logging
    - Health check endpoint
    - Cache layer choredirect hot URLs
    - CI/CD pipeline
    - Horizontal scaling

Trả lời các vấn đề
- concurrency: 2 requests cùng lúc tạo link với cùng URL thì sao?
    - Mỗi request generate short code độc lập
    - MongoDB có unique index trên short_code đảm bảo không có 2 record trùng
    - nếu bị duplicate --> generate + retry

- validation: URL nào hợp lệ (in/output?) + xử lý edges case (unit testing?)
    - hiện tại chỉ check long URL không rỗng
    - các phần sẽ được xử lý nếu có thêm thời gian: validate format URL, giới hạn độ dài input, trả lỗi với message rõ ràng, bổ sung unit test cho các case URL không hợp lệ

- performance: nếu 1 triệu links thì query sao? cần index không?
    - Query theo short_code dùng unique index, độ phức tạp O(log n)
    - Redirect chỉ cần  1 query + 1 update

- security: vấn đề bảo mật cần notice?
    - Hiện dễ bị spam tạo link
    - Chưa có rate limiting
    - Chưa có authentication/authorization
    
- scalability: traffic x100 thì làm gì?
    - CHạy nhiều API instance
    - Db scale bằng sharding/replica set
    - THêm cache
    - Đặt load balancer trước API
- collision: sẽ xảy thường xuyên ra khi số lượng URL tăng tới 1 mức nhất định
    - Check tồn tại trước khi insert
    - unique index trên short_code
    - retry khi có collision

