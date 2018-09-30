import { Model } from '@vuex-orm/core'
import Site from './Site'

export default class Outage extends Model {
  static entity = 'outages'

  static fields() {
    return {
      id: this.attr(null),
      start_time: this.attr(null),
      end_time: this.attr(null),
      description: this.attr(null),
      site_id: this.attr(null),
      site: this.belongsTo(Site, 'site_id')
    }
  }
}
