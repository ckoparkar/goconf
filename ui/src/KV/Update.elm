module KV.Update where

import KV.Model exposing (Model, KV)
import KV.Decoder exposing (decode)

import Effects exposing (Effects, Never)
import Http exposing (get)
import Task exposing (Task)

type Action =
  KVsRetrived (Maybe (List KV))


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
  Http.get decode "http://localhost:8080/v1/kv?recurse=true"
    |> Task.toMaybe
    |> Task.map KVsRetrived
    |> Effects.task
