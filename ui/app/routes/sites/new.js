import Route from '@ember/routing/route';
import RequireAuthMixin from '../../mixins/require-auth-mixin';

export default Route.extend(RequireAuthMixin);
