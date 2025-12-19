# URL shortener (Go)
NỘI DUNG README
Phải có:
- Mô tả bài toán: Bạn hiểu bài toán như thế nào
- Cách chạy project: Từng bước chi tiết (ai cũng chạy được)
- Thiết kế & Quyết định kỹ thuật:
 Tại sao chọn database này?
 Tại sao thiết kế API kiểu này?
 Thuật toán generate mã ngắn là gì?
 Xử lý conflict/duplicate như thế nào?
- Trade-offs:
 "Em chọn X thay vì Y vì lý do A, B, C"
 "Cách này có nhược điểm là... nhưng phù hợp vì..."
- Challenges:
 Gặp vấn đề gì?
 Giải quyết thế nào?
 Học được gì?
- Limitations & Improvements:
 Code hiện tại còn thiếu gì?
 Nếu có thêm thời gian sẽ làm gì?
 Production-ready cần thêm gì?

Trả lời các vấn đề + trong code
- concurrency: 2 requests cùng lúc tạo link với cùng URL thì sao?
- validation: URL nào hợp lệ (in/output?) + xử lý edges case (unit testing?)
- performance: nếu 1 triệu links thì query sao? cần index không?
- security: vấn đề bảo mật cần notice?
- scalability: traffic x100 thì làm gì?
- collision: ở trên nói r

############################################################################


1: define scale (how much URLs need to generate) --> how long should the URLs be
lấy 6 character, total space ~56 tỷ --> ngưỡng nguy hiểm = căn 56 tỷ =238k, tầm này là collision nhiều
--> xài Random base62, check db trước khi insert lại, collision sẽ generate lại

data storage: shortURL=7 bytes, longURL=200 bytes, ...
xài noSQL: 1000 writes need 10k-100k Read per second, able to store big amount of scale, no need joints, complex query
mongoDB vì giải quyết collision/concurrency ằng unique index+retry, phù hợp random base62

API design: REST (simple)
POST
GET

go file: xử lý logic shorten(random base62), trả về original URL, check collision+redo
another way né collision logic: tạo hết URLs, gán vs bool true/false, xài r là true, chưa là false, false thì xài(cái này phải xài SQL: nếu 2 long URL nhập vào, chỉ trỏ tới 1 short URL, trong noSQL lỗi nhưng SQL thì đc vì có acid properties)


ĐỀ: CHO 1 URL, LÀM NÓ NGẮN LẠI, THEO DÕI SỐ LƯỢT CLICK
a) Function
- tạo link rút gọn từ URL dài
- redirect về url gốc
- xem thông tin link (URL gốc, lượt click...)
- liệt kê các links đã tạo

b) Requirements
- xài Golang + database (giải thích lý do chọn)
- code có cấu trúc
- torng README viết rõ cách setup + chạy
- commit Git thường xuyên

c) Trả lời các vấn đề + trong code
- concurrency: 2 requests cùng lúc tạo link với cùng URL thì sao? tạo new shortURL, concurrency k sao
- validation: URL nào hợp lệ (in/output?) + xử lý edges case (unit testing?)
- performance: nếu 1 triệu links thì query sao? cần index không?
- security: vấn đề bảo mật cần notice?
- scalability: traffic x100 thì làm gì?
- collision: ở trên nói r

