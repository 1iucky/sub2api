import landing from './landing'
import common from './common'
import dashboard from './dashboard'
import models from './models'
import admin from './admin'
import misc from './misc'

export default {
  ...landing,
  ...common,
  ...dashboard,
  ...models,
  admin,
  ...misc,
}
