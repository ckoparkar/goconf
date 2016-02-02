import KV.Model exposing (Model)
import KV.View exposing (view)
import KV.Update exposing (init, update)

import StartApp as StartApp
import Task exposing (Task)
import Signal exposing (Signal)
import Html exposing (Html)
import Effects exposing (Never)

app : StartApp.App Model
app =
  StartApp.start
    { init = init
    , update = update
    , view = view
    , inputs = []
    }


main : Signal Html
main =
  app.html


port tasks : Signal (Task.Task Never ())
port tasks =
  app.tasks
