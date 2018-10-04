import Controller from '@ember/controller';

export default Controller.extend({
  auth: Ember.inject.service(),

  actions: {
    login() {
      this.get('auth').login(
        this.get('username'),
        this.get('password')
      ).then(
        (data) => { this.replaceRoute('index'); },
        (jqXHR) => { alert(jqXHR.responseText); }
      );
    }
  }
});
