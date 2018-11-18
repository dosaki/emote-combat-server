package actions

import (
	"github.com/dosaki/emote_combat_server/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/gobuffalo/x/sessions"
	"github.com/rs/cors"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.AllowAll().Handler,
			},
			SessionName: "_emote_combat_server_session",
		})

		// Automatically redirect to SSL
		// app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)

		app.POST("/auth/token", GenerateToken)

		auth := app.Group("/auth")
		auth.Use(RestrictedHandlerMiddleware)

		auth.GET("/invalidate", DestroyToken)

		player := app.Group("/player")
		player.Use(PlayerRestrictedHandlerMiddleware)

		player.GET("/{player_id}", UserList)   // Read
		player.PUT("/{player_id}", UserUpdate) // Update
		app.POST("/player", UsersCreate)       // New

		app.GET("/characters", CharacterList)     // List all
		app.GET("/character/{id}", CharacterList) // Read

		player.GET("/{player_id}/characters", CharacterList)              // Read
		player.GET("/{player_id}/character/{id}", CharacterList)          // Read
		player.POST("/{player_id}/character", CharacterCreate)            // New
		player.PUT("/{player_id}/character/{id}", CharacterUpdate)        // Update
		player.DELETE("/{player_id}/character/{id}", CharacterDelete)     // Delete
		player.GET("/{player_id}/character/{id}/delete", CharacterDelete) // Delete

		app.GET("/character/{character_id}/sheet_entries", SheetEntryList)                        // Read
		player.GET("/{player_id}/character/{character_id}/sheet_entries", SheetEntryList)         // List all
		player.GET("/{player_id}/character/{character_id}/sheet_entry/{id}", SheetEntryList)      // Read
		player.POST("/{player_id}/character/{character_id}/sheet_entry", SheetEntryCreate)        // New
		player.POST("/{player_id}/character/{character_id}/sheet_entries", SheetEntriesCreate)    // New
		player.PUT("/{player_id}/character/{character_id}/sheet_entry/{id}", SheetEntryUpdate)    // Update
		player.PUT("/{player_id}/character/{character_id}/sheet_entries", SheetEntriesUpdate)     // New
		player.DELETE("/{player_id}/character/{character_id}/sheet_entry/{id}", SheetEntryDelete) // Delete

		app.GET("/skills", SkillList)                          // List all
		app.GET("/skill/{id}", SkillList)                      // Read
		app.GET("/skill/{parent_id}/subskills", SkillList)     // Read all subskills
		app.GET("/skill/{parent_id}/subskill/{id}", SkillList) // Read subskill
		app.POST("/skill", SkillCreate)                        // New
		app.PUT("/skill/{id}", SkillUpdate)                    // Update
		app.DELETE("/skill/{id}", SkillDelete)                 // Delete

	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
// func forceSSL() buffalo.MiddlewareFunc {
// 	return forcessl.Middleware(secure.Options{
// 		SSLRedirect:     ENV == "production",
// 		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
// 	})
// }
