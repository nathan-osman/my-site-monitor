import Ember from 'ember';
import Service from '@ember/service';

export default Service.extend({
  user: null,

  /**
   * Attempt to login using the provided credentials
   * @param {String} username
   * @param {String} password
   * @return {Promise}
   */
  login(username, password) {
    return Ember.$.ajax({
      type: 'POST',
      url: '/api/login',
      contentType: 'application/json;charset=utf-8',
      data: JSON.stringify({
        username: username,
        password: password
      })
    }).done((data) => { this.set('user', data); });
  },

  /**
   * Destroy the current session and log the user out
   * @return {Promise}
   */
  logout() {
    return Ember.$.ajax({
      type: 'POST',
      url: '/api/logout'
    }).done(() => { this.set('user', null); });
  }
});
