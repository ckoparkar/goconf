import KV.Model exposing (Model)
import KV.View exposing (view)
import KV.Update exposing (init, update)
import KV.Router

import StartApp as StartApp
import Task exposing (Task)
import Signal exposing (Signal)
import Html exposing (Html)
import Effects exposing (Never)
import RouteHash


app : StartApp.App Model
app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = [
         messages.signal
        ]
    }


main : Signal Html
main =
  app.html


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks


port routeTasks : Signal (Task () ())
port routeTasks =
  RouteHash.start
    { prefix = RouteHash.defaultPrefix
    , address = messages.address
    , models = app.model
    , delta2update = KV.Router.delta2update
    , location2action = KV.Router.location2action
    }


messages : Signal.Mailbox KV.Update.Action
messages =
    Signal.mailbox KV.Update.NoOp
