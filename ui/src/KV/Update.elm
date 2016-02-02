module KV.Update where

import KV.Model exposing (Model, KV)

import Effects exposing (Effects, Never)
import Http exposing (get)
import Json.Decode as Json exposing ((:=))
import Task exposing (Task)

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
