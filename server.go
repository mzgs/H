package H

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/CloudyKit/jet"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/modern-go/reflect2"
	"go.mongodb.org/mongo-driver/bson"
)

var View = jet.NewHTMLSet("./views")

func InitGinServer() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	Server := gin.New()
	Server.Use(gin.Recovery())

	Server.Static("/public", "./public")

	store := cookie.NewStore([]byte("sdvsv#@R@#R$@fvsvsdvdfv"))
	Server.Use(sessions.Sessions("mysession", store))

	return Server

}

func RunGinServer(Server *gin.Engine, PORT, appName string) {
	PL(appName + " start at port: " + PORT)
	Server.Run(":" + PORT)
}

func Render(w io.Writer, view string, variables jet.VarMap, data interface{}) {
	t, err := View.GetTemplate(view + ".html")
	if err != nil {
		PL("GetTemplate ERROR:", err)
	}
	err = t.Execute(w, variables, data)
	if err != nil {
		PL("Execute ERROR:", err)
	}

}

func UploadFile(c *gin.Context, file *multipart.FileHeader) string {

	filename := filepath.Base(file.Filename)
	filename = UrlStringForFile(filename)
	filename = getUniqueFileName(filename, "public/uploads")

	err := c.SaveUploadedFile(file, FM("public/uploads/{name}", filename))
	if err != nil {
		P(err)
	}

	return filename
}

func InitCrudRouters(Server *gin.Engine, DB MongoDBHelper, structsPackageName string) {

	NewFolder("public/uploads")

	typeByName := func(c *gin.Context) reflect2.Type {
		table := c.Param("table")

		return reflect2.TypeByName(structsPackageName + "." + table)
	}

	Server.POST("/crud/list/:table", func(c *gin.Context) {

		t := typeByName(c)

		slice := reflect.MakeSlice(reflect.SliceOf(t.Type1()), 0, 0)
		x := reflect.New(slice.Type())
		x.Elem().Set(slice)

		jtPageSize := Int(c.Query("jtPageSize"))
		jtStartIndex := Int(c.Query("jtStartIndex"))
		jtSorting := c.Query("jtSorting")

		filters := bson.M{}
		_ = filters

		c.Request.ParseForm()
		for k, v := range c.Request.PostForm {
			PL(k, v)
		}

		dbQuery := DB.Find(x.Interface())

		if jtSorting != "" {
			splitComma := strings.Split(jtSorting, ",")

			var sortArray []string
			for _, c := range splitComma {
				split := strings.Split(c, " ")
				key := strings.ToLower(split[0])
				direction := split[1]

				if direction == "DESC" {
					key = "-" + key
				}
				sortArray = append(sortArray, key)
			}

			dbQuery = dbQuery.Sort(sortArray...)

		}

		count := DB.Count(x.Interface(), bson.M{})
		dbQuery.Limit(jtPageSize).Skip(jtStartIndex).All()

		rows := x.Elem().Interface()

		c.JSON(200, gin.H{
			"Result":           "OK",
			"TotalRecordCount": count,
			"Records":          rows,
		})

	})

	Server.POST("/crud/create/:table", func(c *gin.Context) {
		table := c.Param("table")
		t := reflect2.TypeByName("main." + table)

		x := reflect.New(t.Type1())
		c.Bind(x.Interface())

		record := x.Elem().Interface()
		DB.InsertOne(record)

		c.JSON(200, gin.H{
			"Result": "OK",
			"Record": record,
		})

	})

	Server.POST("/crud/edit/:table", func(c *gin.Context) {
		table := c.Param("table")
		t := reflect2.TypeByName("main." + table)

		x := reflect.New(t.Type1())
		c.Bind(x.Interface())

		id := reflect.ValueOf(DB.ObjectId(c.PostForm("_id")))
		x.Elem().FieldByName("ID").Set(id)

		record := x.Elem().Interface()

		DB.UpdateOne(record)

		c.JSON(200, gin.H{
			"Result": "OK",
			"Record": record,
		})

	})

	Server.POST("/crud/delete/:table", func(c *gin.Context) {

		table := c.Param("table")
		t := reflect2.TypeByName("main." + table)

		err := DB.DeleteOne(reflect.New(t.Type1()).Elem().Interface(), bson.M{"_id": DB.ObjectId(c.PostForm("_id"))})

		if err != nil {
			c.JSON(200, gin.H{
				"Result":  "ERROR",
				"Message": err.Error(),
			})

			return
		}

		c.JSON(200, gin.H{
			"Result": "OK",
		})

	})

	Server.POST("/image", func(c *gin.Context) {

		m, _ := c.MultipartForm()

		for fieldName, _ := range m.File {

			file, err := c.FormFile(fieldName)
			if err != nil {
				PL("image upload error /image", err)
			}
			filename := UploadFile(c, file)

			// tinymce file upload
			if fieldName == "file" {

				c.JSON(200, gin.H{
					"location": filename,
				})
				return
			}

			c.String(200, filename)
			return

		}

	})
	Server.GET("/image", func(c *gin.Context) {
		load := c.Query("load")

		//c.String(200,  "uploads/" + load)
		c.File("public/uploads/" + load)

	})

	Server.GET("/uploads", func(c *gin.Context) {

		files, _ := ioutil.ReadDir("public/uploads")

		imageExt := []string{".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp", ".svg", ".pdf"}

		var images []string
		for _, k := range files {
			if !Contains(imageExt, filepath.Ext(k.Name())) {
				continue
			}

			images = append(images, k.Name())
		}

		c.JSON(200, images)

	})

	Server.GET("/resimler", func(c *gin.Context) {

		files, _ := ioutil.ReadDir("public/uploads")

		sort.Slice(files, func(i, j int) bool {
			return files[i].ModTime().Before(files[j].ModTime())
		})

		imageExt := []string{".png", ".jpg", ".jpeg", ".gif", ".ico", ".webp", ".svg", ".pdf"}

		var images []string
		for _, k := range files {
			if !Contains(imageExt, filepath.Ext(k.Name())) {
				continue
			}

			images = append(images, k.Name())
		}

		v := jet.VarMap{}
		v.Set("title", "Resimler")

		Render(c.Writer, "resimler", v, images)

	})

	Server.POST("/delete-image", func(c *gin.Context) {
		img := c.PostForm("img")
		os.Remove("public/uploads/" + img)
		c.JSON(200, true)
	})

	Server.GET("/downloadpng", func(c *gin.Context) {
		load := c.Query("load")

		outFile := "public/uploads/" + load

		c.FileAttachment(outFile, filepath.Base(outFile))
	})

}
