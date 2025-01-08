import {Routes} from '@angular/router';
import {DeploymentsComponent} from '../components/deployments/deployments.component';
import {PodsComponent} from '../components/pods/pods.component';

export const routes: Routes = [
  {
    path: '',
    children: [
      {
        path: 'deployments',
        component: DeploymentsComponent
      },
      {
        path: 'pods',
        component: PodsComponent
      }
    ],
  },

];
