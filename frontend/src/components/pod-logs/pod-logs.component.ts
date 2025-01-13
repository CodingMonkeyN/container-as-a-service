import {Component, HostBinding, inject, Signal} from '@angular/core';
import {ApiService} from '../../services/api.service';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {takeUntilDestroyed, toSignal} from '@angular/core/rxjs-interop';
import {
  MatCard,
  MatCardActions,
  MatCardContent, MatCardFooter,
  MatCardHeader,
  MatCardSubtitle,
  MatCardTitle
} from '@angular/material/card';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatIcon} from '@angular/material/icon';
import {BehaviorSubject, switchMap} from 'rxjs';

@Component({
  selector: 'app-pod-logs',
  imports: [
    MatCard,
    MatCardContent,
    MatCardHeader,
    MatCardTitle,
    MatCardSubtitle,
    MatCardActions,
    RouterLink,
    MatButton,
    MatCardFooter,
    MatIconButton,
    MatIcon,
  ],
  templateUrl: './pod-logs.component.html',
})
export class PodLogsComponent {

  @HostBinding() class = 'flex-1 p-4 flex flex-col h-full w-full';
  logs: Signal<string>
  namespace: string;
  containerName: string;
  refresh = new BehaviorSubject(undefined)

  constructor(route: ActivatedRoute, api: ApiService) {
    console.log(route.snapshot.params)
    this.namespace = route.snapshot.params['namespace'];
    this.containerName = route.snapshot.params['containerName'];

    this.logs = toSignal(this.refresh.pipe(
        takeUntilDestroyed(),
        switchMap(() => api.getLogs(this.namespace, this.containerName)))
      , {initialValue: ''});
  }

  back() {
    window.history.back();
  }
}
