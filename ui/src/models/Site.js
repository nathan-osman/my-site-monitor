import { Model } from '@vuex-orm/core'
import User from './User'

export default class Site extends Model {
  static entity = 'sites'

  static fields() {
    return {
      id: this.attr(null),
      url: this.attr(null),
      name: this.attr(null),
      poll_interval: this.attr(null),
      poll_time: this.attr(null),
      status: this.attr(null),
      status_time: this.attr(null),
      user_id: this.attr(null),
      user: this.belongsTo(User, 'user_id')
    }
  }
}
