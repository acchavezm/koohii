# Proyecto de Go 
Proyecto de Go realizado para la materia de Principios de Lenguajes de Programación.

## Descripción
La aplicación realiza recomendaciones de las pistas (tracks) del Top 50 - Global de Spotify en base al clima (puntualmente, la temperatura) de la Ciudad de Guayaquil. El valor de "energía" de los tracks es utilizado para esta simple comparación. Tracks con alta energía se consideran más afines a temperaturas altas, mientras que temperaturas bajas a los de baja energía.

En un principio, se solicitan las credenciales de Spotify, manejando una sesión temporal utilizando el token de autenticación proveído por [el flujo de autenticación de Spotify](https://developer.spotify.com/documentation/general/guides/authorization-guide/#authorization-code-flow). Se carga el nombre, foto de perfil y contador de seguidores del usuario que iniciara sesión. Se muestra también la temperatura de la ciudad de Guayaquil y la descripción del clima. También, en base a esta data, se muestran tracks recomendados para esas condiciones, tal como se describió anteriormente.
## Acerca de la implementación
### Frameworks utilizados
- [Gin Web Framework](https://github.com/gin-gonic/gin): utlilizado para el levantamiento del servidor, routing, statics y templates.
- [Bootstrap 5.1](https://getbootstrap.com/docs/5.1/getting-started/introduction/): principalmente para UI.
- [Biblioteca de Spotify para Go](https://github.com/zmb3/spotify): wrapper para realizar consultas al Web API de Spotify.
- [Favicon Gin's Middleware](https://github.com/thinkerou/favicon): para el favicon.

## APIs utilizados
- OpenWeather - Current Weather Data. La documentación se encuentra [aquí](https://openweathermap.org/current). El endpoint utilizado fue el GET de current weather data con el nombre de la ciudad (Guayaquil).
- Spotify - Web API. La referencia se encuentra [aquí](https://developer.spotify.com/documentation/web-api/reference/). Se usaron los APIs de User Profile (fetch de info.), Playlists (GET tracks del playlist), Tracks (GET audio features de una lista de IDs).

## Instalación y uso
1. Ir al [portal de desarrolladores de Spotify](https://developer.spotify.com/dashboard/login), iniciar sesión y crear una aplicación con su nombre y descripción.
2. Una vez en el dashboard, copiar los strings de **Client ID** y **Client Secret**. Estos serán utilizados por la aplicación para obtener un token de autenticación, y en base a aquello realizar las llamadas al Web API.
3. Darle clic a **Edit Settings**, y definir la URI de redirección. En este caso, definir http://localhost:9001/callback. Se puede usar otro puerto en lugar de 9001, verificar que el cambio también se haga en main.go. Tener copiado este valor también.
4. Darle clic a **Save**.
5. Esto es todo por Spotify. Seguimos con OpenWeather.
6. Ingresar a https://home.openweathermap.org/users/sign_up. Crear la cuenta para poder obtener el API Key.
7. Una vez loggeado, ir a https://home.openweathermap.org/api_keys. Copiar el valor de **Key**.
8. Crear un archivo .env en la raíz del proyecto, el archivo debe de contener la siguiente estructura:


        SPOTIFY_ID=<El Client ID de Spotify>
        SPOTIFY_SECRET=<El Client Secret de Spotify>
        SPOTIFY_REDIRECT_URI=<El Redirect URI definido en Spotify>
        WEATHER_API_KEY=<El API Key proveído por OpenWeather>

9. Esto es todo para el tema de configuración de los servicios.
10. En el terminal darle a go run main.go. Acceder a http://localhost:9001/index.