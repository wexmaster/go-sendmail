# go-sendmail 

Envió de correos con GO usando SMTP y credenciales. La intención de esta aplicacion en GO es el uso simple para enviar desde Windows o Linux correos con formato HTML y imágenes en BASE64.


Es necesario tener instalado GO en el equipo que necesite compilar. Sea Windows o Linux.

dentro del directorio donde este main.go
Ejecute:
```go mod init```
Luego:
```go mod tidy```
Y por ultimo compilamos:
```go build main.go```

Tendras un archivo main o main.exe

Para mas ayuda en la consula ejecute main.exe --help

```./main --smtp=mail.SUHOST.es --puerto=25 --from=SUCORREO@HOST.es --password='CLAVE' --sendto=CORREODESSTINATARIO@hotmail.com```

Recibirá un correo por defecto:

![Alt Text](https://wexmaster.es/img/ejemplo.png)
