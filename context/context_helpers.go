package context

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
)

// bind helper function
// Hàm nội bộ để liên kết các giá trị từ url.Values vào một struct.
// Sử dụng reflection để map các giá trị vào các trường struct dựa trên tag "form" hoặc "json".
//
// Parameters:
//   - values: Các giá trị cần được liên kết vào struct
//   - obj: Con trỏ đến struct sẽ nhận các giá trị
//
// Returns:
//   - error: Lỗi nếu không thể liên kết giá trị
//
// Errors:
//   - "obj must be a non-nil pointer": Khi đối tượng không phải là con trỏ hoặc là nil
//   - "obj must be a struct": Khi đối tượng không phải là struct
func bind(values url.Values, obj interface{}) error {
	// Kiểm tra xem đối tượng có phải là con trỏ không null hay không
	objValue := reflect.ValueOf(obj)
	if objValue.Kind() != reflect.Ptr || objValue.IsNil() {
		return errors.New("obj must be a non-nil pointer")
	}

	// Lấy giá trị thực của đối tượng
	objValue = objValue.Elem()
	objType := objValue.Type()

	// Kiểm tra xem đối tượng có phải là struct hay không
	if objType.Kind() != reflect.Struct {
		return errors.New("obj must be a struct")
	}

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		formTag := field.Tag.Get("form")
		if formTag == "" {
			formTag = field.Tag.Get("json") // Fallback to json tag
		}
		if formTag == "" || formTag == "-" {
			continue
		}

		formValue := values.Get(formTag)
		if formValue == "" {
			continue
		}

		fieldValue := objValue.Field(i)
		if !fieldValue.CanSet() {
			continue
		}

		err := setFieldValue(fieldValue, formValue)
		if err != nil {
			return fmt.Errorf("binding error for field %s: %w", field.Name, err)
		}
	}

	return nil
}

// setFieldValue đặt giá trị cho trường dựa trên đầu vào chuỗi.
// Hàm này chuyển đổi giá trị chuỗi thành kiểu dữ liệu tương ứng của trường
// và gán giá trị đã chuyển đổi vào trường đó sử dụng reflection.
//
// Parameters:
//   - fieldValue: Giá trị trường cần đặt (reflect.Value)
//   - value: Giá trị chuỗi cần chuyển đổi và gán
//
// Returns:
//   - error: Lỗi nếu có trong quá trình chuyển đổi kiểu
//
// Errors:
//   - strconv: Lỗi chuyển đổi chuỗi sang kiểu số
//   - "unsupported field type": Kiểu dữ liệu không được hỗ trợ
func setFieldValue(fieldValue reflect.Value, value string) error {
	// Xử lý tùy theo kiểu dữ liệu của trường
	switch fieldValue.Kind() {
	case reflect.String:
		// Đối với chuỗi, gán trực tiếp
		fieldValue.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Đối với các kiểu số nguyên có dấu
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Đối với các kiểu số nguyên không dấu
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		fieldValue.SetUint(val)
	case reflect.Float32, reflect.Float64:
		// Đối với các kiểu số thực
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fieldValue.SetFloat(val)
	case reflect.Bool:
		// Đối với kiểu boolean
		val, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		fieldValue.SetBool(val)
	default:
		// Trả về lỗi cho các kiểu không được hỗ trợ
		return fmt.Errorf("unsupported field type: %s", fieldValue.Kind())
	}
	return nil
}
