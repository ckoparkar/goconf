module KV.View where

import KV.Model exposing (Model, KV, uniqueKVs)
import KV.Update exposing (Action)
import KV.Utils exposing (kvsForActivePage)

import Html exposing (..)
import Html.Attributes exposing (..)
import Signal exposing (Signal, Address)


kvRow : String -> KV -> Html
kvRow activePage kv =
  tr [] [ td [] [a [href ("#!/" ++ activePage ++ kv.key)] [text kv.key] ]
        , td [] [text kv.value]
        ]


view : Signal.Address Action -> Model -> Html
view address model =
  div [] [ table [] [
                    tbody []
                      ( kvsForActivePage model
                      |> uniqueKVs
                      |> List.map (kvRow model.activePage))
                   ] ]
