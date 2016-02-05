module KV.Router where

import KV.Model exposing (Model)
import KV.Update exposing (Action)

import RouteHash exposing (HashUpdate)
import String


delta2update : Model -> Model -> Maybe HashUpdate
delta2update previous current =
  Just <| RouteHash.set (current.activePage :: [])


location2action : List String -> List Action
location2action list =
  ( KV.Update.SetActivePage (String.join "/" list) ) :: []
