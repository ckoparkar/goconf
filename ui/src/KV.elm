module KV where

import Effects exposing (..)
import StartApp as StartApp
import Task exposing (Task)
import Html exposing (..)
import Html.Attributes exposing (..)
import Signal exposing (Signal, Address)
import Http exposing (..)
import Json.Decode as Json exposing ((:=))

-- model

type alias KV =
  { key : String
  , value : String
  }


type alias Model =
  { kvs : List KV
  }

-- update

type Action
  = KVsRetrived (Maybe (List KV))


init : (Model, Effects Action)
init =
  ( Model []
    , getAllKVs
  )


update : Action -> Model -> (Model, Effects Action)
update action model =
    case action of
      KVsRetrived xs ->
        ( { model | kvs = (Maybe.withDefault [] xs) }
        , Effects.none
        )

getAllKVs : Effects.Effects Action
getAllKVs =
  Http.get kvs "http://localhost:8080/v1/kv?recurse=true"
    |> Task.toMaybe
    |> Task.map KVsRetrived
    |> Effects.task


kv : Json.Decoder KV
kv =
  Json.object2 KV
    ("key" := Json.string)
    ("value" := Json.string)


kvs : Json.Decoder (List KV)
kvs =
  Json.list kv


-- view

kvRow : KV -> Html
kvRow kv =
  tr [] [
     td [] [text (toString kv.key)]
    ,td [] [text kv.value]
  ]


view : Signal.Address Action -> Model -> Html
view address model =
  div [class "container-fluid"] [
        h1 [] [text "KVs" ]
      , table [class "table table-striped"] [
          thead [] [
            tr [] [
               th [] [text "key"]
              ,th [] [text "Value"]
          ]
        ]
      , tbody [] (List.map kvRow model.kvs)
    ]
  ]


-- Main

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
