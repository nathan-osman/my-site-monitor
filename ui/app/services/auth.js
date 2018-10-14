import Ember from 'ember';
import Service from '@ember/service';

export default Service.extend({
  user: null,

  /**
   * Attempt to restore a session by obtaining the user's info
   * @return {Promise}
   */
  restore() {
    return new Ember.RSVP.Promise((resolve) => {
      Ember.$.ajax({
        url: '/api/users/me'
      })
      .done((data) => { this.set('user', data); })
      .always(() => { resolve(); })
    });
  },

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
    })
    .done((data) => { this.set('user', data); });
  },

  /**
   * Destroy the current session and log the user out
   * @return {Promise}
   */
  logout() {
    return Ember.$.ajax({
      type: 'POST',
      url: '/api/logout'
    })
    .done(() => { this.set('user', null); });
  }
});
