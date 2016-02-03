module KV.View where

import KV.Model exposing (Model, KV, uniqueKVs)
import KV.Update exposing (Action)

import Html exposing (..)
import Html.Attributes exposing (class)
import Signal exposing (Signal, Address)

kvRow : KV -> Html
kvRow kv =
  tr [] [ td [] [text kv.key]
        , td [] [text kv.value]
        ]


view : Signal.Address Action -> Model -> Html
view address model =
  div [] [ table [] [ tbody [] (List.map kvRow (uniqueKVs model.kvs)) ] ]
