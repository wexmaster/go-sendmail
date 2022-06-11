package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/toorop/go-dkim"
	mail "github.com/xhit/go-simple-mail/v2"
)

const privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAwGrscUxi9zEa9oMOJbS0kLVHZXNIW+EBjY7KFWIZSxuGAils
wBVl+s5mMRrR5VlkyLQNulAdNemg6OSeB0R+2+8/lkHiMrimqQckZ5ig8slBoZhZ
wUoL/ZkeQa1bacbdww5TuWkiVPD9kooT/+TZW1P/ugd6oYjpOI56ZjsXzJw5pz7r
DiwcIJJaaDIqvvc5C4iW94GZjwtmP5pxhvBZ5D6Uzmh7Okvi6z4QCKzdJQLdVmC0
CMiFeh2FwqMkVpjZhNt3vtCo7Z51kwHVscel6vl51iQFq/laEzgzAWOUQ+ZEoQpL
uTaUiYzzNyEdGEzZ2CjMMoO8RgtXnUo2qX2FDQIDAQABAoIBAHWKW3kycloSMyhX
EnNSGeMz+bMtYwxNPMeebC/3xv+shoYXjAkiiTNWlfJ1MbbqjrhT1Pb1LYLbfqIF
1csWum/bjHpbMLRPO++RH1nxUJA/BMqT6HA8rWpy+JqiLW9GPf2DaP2gDYrZ0+yK
UIFG6MfzXgnju7OlkOItlvOQMY+Y501u/h6xnN2yTeRqXXJ1YlWFPRIeFdS6UOtL
J2wSxRVdymHbGwf+D7zet7ngMPwFBsbEN/83KGLRjkt8+dMQeUeob+nslsQofCZx
iokIAvByTugmqrB4JqhNkAlZhC0mqkRQh7zUFrxSj5UppMWlxLH+gPFZHKAsUJE5
mqmylcECgYEA8I/f90cpF10uH4NPBCR4+eXq1PzYoD+NdXykN65bJTEDZVEy8rBO
phXRNfw030sc3R0waQaZVhFuSgshhRuryfG9c1FP6tQhqi/jiEj9IfCW7zN9V/P2
r16pGjLuCK4SyxUC8H58Q9I0X2CQqFamtkLXC6Ogy86rZfIc8GcvZ9UCgYEAzMQZ
WAiLhRF2MEmMhKL+G3jm20r+dOzPYkfGxhIryluOXhuUhnxZWL8UZfiEqP5zH7Li
NeJvLz4pOL45rLw44qiNu6sHN0JNaKYvwNch1wPT/3/eDNZKKePqbAG4iamhjLy5
gjO1KgA5FBbcNN3R6fuJAg1e4QJCOuo55eW6vFkCgYEA7UBIV72D5joM8iFzvZcn
BPdfqh2QnELxhaye3ReFZuG3AqaZg8akWqLryb1qe8q9tclC5GIQulTInBfsQDXx
MGLNQL0x/1ylsw417kRl+qIoidMTTLocUgse5erS3haoDEg1tPBaKB1Zb7NyF8QV
+W1kX2NKg5bZbdrh9asekt0CgYA6tUam7NxDrLv8IDo/lRPSAJn/6cKG95aGERo2
k+MmQ5XP+Yxd+q0LOs24ZsZyRXHwdrNQy7khDGt5L2EN23Fb2wO3+NM6zrGu/WbX
nVbAdQKFUL3zZEUjOYtuqBemsJH27e0qHXUls6ap0dwU9DxJH6sqgXbggGtIxPsQ
pQsjEQKBgQC9gAqAj+ZtMXNG9exVPT8I15reox9kwxGuvJrRu/5eSi6jLR9z3x9P
2FrgxQ+GCB2ypoOUcliXrKesdSbolUilA8XQn/M113Lg8oA3gJXbAKqbTR/EgfUU
kvYaR/rTFnivF4SL/P4k/gABQoJuFUtSKdouELqefXB+e94g/G++Bg==
-----END RSA PRIVATE KEY-----`

func main() {
	server := mail.NewSMTPClient()
	// Pasar una imagen PNG a Base64 https://onlinepngtools.com/convert-png-to-base64
	htmlBody := flag.String("htmlBody", "<html><head><meta http-equiv='Content-Type' content='text/html; charset=utf-8' /><title>Hola Amig@!</title></head><body><p>Este es un programa <b>hecho en GO por wexmaster para compilar en Windows o Linux</b>.</p><p><img src='data:image/png;base64,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' alt='https://github.com/wexmaster' /></p><p>By wexmaster</p></body></html>", "a string")

	smtp := flag.String("smtp", "127.0.0.1", "a string")
	puerto := flag.Int("puerto", 25, "an int")
	name := flag.String("name", "Admin", "a string")
	from := flag.String("from", "admin@localhost.es", "a string")
	title := flag.String("title", "Correo de ejemplo.", "a string")
	password := flag.String("password", "", "a string")
	// Destinatatio
	sendto := flag.String("sendto", "QuienRecive@gmail.com", "a string")

	// Usado para imprimir ayuda del programa y parametros
	help := flag.Bool("help", false, "a bool")
	// Usado para activar alta prioridad
	altaprioridad := flag.Bool("HighPriority", false, "a bool")
	// Usado para adjuntar un archivo
	attach := flag.String("Attach", "none", "a string")

	flag.Parse()

	if *help {

		fmt.Println("Ayuda para enviar un correo!!")
		fmt.Println("")
		fmt.Println(" EJEMPLO : sendmail.exe --smtp=smtp.gmail.com --puerto=587 --from=QuienEnvia@gmail.com --password='CLAVE-QuienEnvia' --sendto=QuienRecive@gmail.com")
		fmt.Println("")
		fmt.Println("    Parámetro --smtp=HOST-IP Es necesario especificar el servidor por defecto esta 127.0.0.1")
		fmt.Println("    Parámetro --puerto=25 Por defecto es 25 pero puede ser 587 o otro a definir en su servidor.")
		fmt.Println("    Parámetro --name='Tu Nombre' Por defecto es Admin.")
		fmt.Println("    Parámetro --title='Titulo del correo o Asunto' Por defecto es 'Correo de ejemplo.")
		fmt.Println("    Parámetro --from=QuienEnvia@gmail.com Por defecto admin@localhost.es")
		fmt.Println("    Parámetro --password='CLAVEdeCorreo' Por defecto esta vacio.")
		fmt.Println("    Parámetro --sendto=QuienRecive@gmail.com Por defecto el correo es QuienRecive@gmail.com")
		fmt.Println("    Parámetro --altaprioridad Activa la alta prioridad en el mensaje")
		fmt.Println("    Parámetro --htmlBody='<p>Hola <strong>Mundo</strong>!</p>' Aquí puede usar HTML o texto en una línea para enviarlo en el cuerpo de mensaje.")
		fmt.Println("    Parámetro --Attach='C:/Users/Users/Desktop/logo.png' Para adjuntar un archivo colocar la ruta del mismo.")
		fmt.Println("")
		fmt.Println(" Si usa GMAIL verificar tener activado IMAP. Configuración -> Ver todo los Ajustes -> Reenvió y Correo POP/IMAP.")
		fmt.Println("    Otro parámetro de GMAIL es desactivar para que otras Aplicaciones menos seguras puedan enviar correo:")
		fmt.Println("    https://www.google.com/settings/security/lesssecureapps")
		fmt.Println("  !!!!!!!! GOOGLE TIENE DESACTIVADO ESTE METODO ¡¡¡¡¡¡¡¡¡")
		fmt.Println("")

		os.Exit(0)

	}

	// Datos de SMTP
	server.Host = *smtp
	server.Port = *puerto
	server.Username = *from
	server.Password = *password
	server.Encryption = mail.EncryptionSTARTTLS

	setfrom := fmt.Sprint(*name + " <" + *from + ">")
	addto := fmt.Sprint(*sendto)
	subjet := fmt.Sprint(*title)

	server.Authentication = mail.AuthPlain

	server.KeepAlive = false

	// Timeout de conexion con el SMTP
	server.ConnectTimeout = 10 * time.Second

	// Timeout para enviar los datos y esperar respuesta
	server.SendTimeout = 10 * time.Second

	// Si es verdadero no verificara el certificado tls si usa web internas mejor habilitado para que no falle
	server.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// cliente SMTP
	smtpClient, err := server.Connect()

	if err != nil {
		log.Fatal(err)
	}

	// Nuevo correo con HTML
	email := mail.NewMSG()
	email.SetFrom(setfrom).
		AddTo(addto).
		SetSubject(subjet)

	email.SetBody(mail.TextHTML, *htmlBody)
	//Envia con Prioridad si se adjunta con parametros
	if *altaprioridad {
		email.SetPriority(mail.PriorityHigh)
	}

	// Añade imagen dentro del HTML para usarlo en la firma por ejemplo
	// email.Attach(&mail.File{FilePath: "C:/Users/Usuario/Desktop/logo.png", Name: "logo.png", Inline: true})
	// Añade imagen o archivo como adjunto en el correo para enviar PDF o ZIP
	// email.Attach(&mail.File{FilePath: "C:/Users/User/Desktop/logo.png"})
	// email.Attach(&mail.File{FilePath: "C:/Users/User/Pictures/Gopher.png", Inline: true})

	if *attach != "none" {
		email.Attach(&mail.File{FilePath: *attach})
	}

	// https://www.dmarcanalyzer.com/es/dkim-3/
	// Usado pero no se a creado oficialmente una KEY valida por lo que usaremos el del propietario de la libreria
	if privateKey != "" {
		options := dkim.NewSigOptions()
		options.PrivateKey = []byte(privateKey)
		options.Domain = "example.com"
		options.Selector = "default"
		options.SignatureExpireIn = 3600
		options.Headers = []string{"from", "date", "mime-version", "received", "received"}
		options.AddSignatureTimestamp = true
		options.Canonicalization = "relaxed/relaxed"

		email.SetDkim(options)
	}

	// Verifica si hay error despues del envio
	if email.Error != nil {
		log.Fatal(email.Error)
	}

	// LLama a Send y pasa al cliente
	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Correo enviado!!")
	}
}
