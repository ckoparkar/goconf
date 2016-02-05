module KV.Utils where

import KV.Model exposing (Model, KV)

import Regex exposing (regex, replace)
import String

kvsForActivePage : Model -> List KV
kvsForActivePage model =
  List.filter (\kv -> String.startsWith model.activePage kv.key) model.kvs
    |> List.map (\kv -> KV (replace Regex.All (regex model.activePage) (\_ -> "") kv.key) kv.value)
