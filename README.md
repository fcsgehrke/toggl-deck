# toggl-deck
## Toggl Backend Unattended Programming Test

That's a test for the of Backend Software Engineer at Toggl. It's based on their description [here](https://toggl.notion.site/Toggl-Backend-Unattended-Programming-Test-015a95428b044b4398ba62ccc72a007e).

### Starting

To start this app you're gonna need Docker with Compose and Make from Linux or Mac. 
Once you have it installed. It's just to type:

> make up

It'll download all the dependencies and run the app on http://localhost:3000. There's a swagger UI on http://localhost:3000/docs/index.html. You can see the docs there and try the APIs too in the UI.

### Layout
I tried to follow the [Golang Standard Project Layout](https://github.com/golang-standards/project-layout) for this project. Usually in all the companies I worked they always have their own layout but always inspired by the Stardard one.

### Libraries
The app is using most of the standard libraries for testing, logging and others. I just added [chi](https://github.com/go-chi/chi) as a router for the HTTP api. A [UUID library](https://github.com/gofrs/uuid). [GORM](gorm.io) as the ORM. Postgres as the DB. There's also pg-admin which makes the development easier. I'm sure you know all of it.

Well, I'm not gonna tell you all the details here. I'll let some fun to the code review.

### Possible Improvements
- Use a better log library like logrus
- Use a better test library (we're using go standard test library) like testify
- Add hot-reloading in the docker build with comstrek/air
- Add OpenTelemetry with support for chi, gorm
- Add some kind of auth 
- Add CI/CD for k8s or Cloud Functions
