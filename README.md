# ProjectGo
url API:http://localhost:3000/recommended
Body của api có định dạng như sau:
Id 10
latitude  5.108851
longitude 97.578896
distance  10.0
ageStart  10
ageEnd  40
gender  M
ignore_array  [1,4,10,19]
Giải pháp thực hiện để đảm bảo data recommended trả về không chứa users đã bị ignore khi độ lớn của ignore_aray này tăng lên
-> Giải pháp:
Cách 1: Sử dụng kỹ thuật worker & go routines trong golang
- Tạo ra các go routines bằng cách ta chọn mặc định số lượng worker thực hiện: JobsChannel lưu trữ khối lượng công việc,Results lưu trữ kết quả trả về
- Tham số truyền vào cho mỗi job trong channel jobs là các subArray được phân chia ra từ mảng user ban đầu bằng cách chia đều só lượng phần của của mảng user cho số lượng worker
,tham sô thứ hai là ignore_aray
- Mỗi woker sẽ tạo ra một luồng xử lý dựa trên core phân bổ cho nó. Trong luồng xứ lý này sẽ thực hiện lọc và loại bỏ các user có id ở trong mảng user
Cách 2: Dựa vào thời gian cần đến của bộ dữ liệu để ta phân bổ thành nhiều lần gửi đến client
- Chia nhỏ khối lượng mảng user ban đầu thành nhiều phần 
- Lọc và trả về dữ liệu cần ứng với mỗi phần chia nhỏ. Server gửi dữ liệu dần dần theo thời gian
