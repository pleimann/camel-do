/* @refresh reload */
import {render} from 'solid-js/web'
import {OauthService} from "@bindings/pleimann.com/camel-do/services/oauth";
import App from './App'

OauthService.Google()
  .then(() => {
    render(() => <App/>, document.body!)
  })

