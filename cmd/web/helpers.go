package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.sessionManager.Exists(r.Context(), "authenticatedUserID")
}

// Create an newTemplateData() helper, which returns a pointer to a templateData
// struct initialized with the current year. Note that we're not using the
// *http.Request parameter here at the moment, but we will do later in the book.
// 这里返回一个初始化好当前年份的 templateData 结构
func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
		// Add the flash message to the template data, if one exists.
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	err = app.errorLog.Output(2, trace)
	if err != nil {
		return
	}
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initialize a new buffer.
	// 初始化新 buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError() helper
	// and then return.
	// 将拼接好的模板写入 buffer，而不是直接调用 http.ResponseWriter， 如果有报错，调用 serverError() 辅助函数并返回
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the template is written to the buffer without any errors, we are safe
	// to go ahead and write the HTTP status code to http.ResponseWriter.
	// 如果模板拼接没有错误，给响应设置正确的 HTTP 状态码
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter. Note: this
	// is another time where we pass our http.ResponseWriter to a function that
	// takes an io.Writer.
	// 将缓冲区内容写入 http.ResponseWriter
	_, err = buf.WriteTo(w)
	if err != nil {
		return
	}
}

// Create a new decodePostForm() helper method. The second parameter here, dst,
// is the target destination that we want to decode the form data into.
// 创建一个新的 decodePostForm() 帮助函数，第二个参数 dst 表示 form 数据传递到哪里
func (app *application) decodePostForm(r *http.Request, dst any) error {
	// Call ParseForm() on the request, in the same way that we did in our
	// createSnippetPost handler.
	// 调用 ParseForm()
	err := r.ParseForm()
	if err != nil {
		return err
	}

	// Call Decode() on our decoder instance, passing the target destination as
	// the first parameter.
	// 调用 Decode() 方法，dst 作为第一个参数
	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// If we try to use an invalid target destination, the Decode() method
		// will return an error with the type *form.InvalidDecoderError.We use
		// errors.As() to check for this and raise a panic rather than returning
		// the error.
		// 如果使用了错误的 dst， Decode() 方法会返回一个类型为 *form.InvalidDecoderError 的错误，
		// 使用 errors.As() 检查指定的错误类型，并引发一个 panic，这样处理比直接返回错误更好。
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// For all other errors, we return them as normal.
		// 如果是其他错误，那么会返回错误
		return err
	}
	// 没有错误发生会返回 nil，供上层代码捕捉错误
	return nil
}
