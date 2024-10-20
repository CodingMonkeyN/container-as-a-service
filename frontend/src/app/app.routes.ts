import {Routes} from '@angular/router';
import {DeploymentsComponent} from '../components/deployments/deployments.component';
import {PodsComponent} from '../components/pods/pods.component';
import {PodLogsComponent} from '../components/pod-logs/pod-logs.component';

export const routes: Routes = [
  // Umleitung von leerem Pfad auf 'deployments'
  {
    path: '',
    redirectTo: 'deployments',
    pathMatch: 'full', // Nur bei vollständigem Leerpfad umleiten
  },
  {
    path: 'deployments',
    component: DeploymentsComponent
  },
  {
    path: 'pods',
    children: [
      {
        path: '',
        component: PodsComponent,
      },
      {
        path: 'logs/:namespace/:containerName',
        component: PodLogsComponent
      }
    ]
  }
];
