import {Component, HostBinding, inject, Signal} from '@angular/core';
import {
  MatCell,
  MatCellDef,
  MatColumnDef,
  MatHeaderCell, MatHeaderCellDef,
  MatHeaderRow,
  MatHeaderRowDef,
  MatRow, MatRowDef, MatTable
} from "@angular/material/table";
import {toSignal} from '@angular/core/rxjs-interop';
import {ApiService} from '../../services/api.service';
import {Pod} from '../../models/pod.model';
import {MatButton, MatIconButton} from '@angular/material/button';
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from '@angular/material/card';
import {BehaviorSubject, switchMap} from 'rxjs';
import {MatIcon} from '@angular/material/icon';
import {Router, RouterLink} from '@angular/router';

@Component({
  selector: 'app-pods',
  imports: [
    MatCell,
    MatCellDef,
    MatColumnDef,
    MatHeaderCell,
    MatHeaderRow,
    MatHeaderRowDef,
    MatRow,
    MatRowDef,
    MatTable,
    MatHeaderCellDef,
    MatButton,
    MatCard,
    MatCardHeader,
    MatCardTitle,
    MatCardContent,
    MatIcon,
    MatIconButton,
    RouterLink
  ],
  templateUrl: './pods.component.html',
})
export class PodsComponent {
  @HostBinding() class = 'flex-1 flex flex-col h-full w-full';

  #refresh = new BehaviorSubject<void>(undefined)
  displayedColumns: string[] = ['namespace', 'name', 'cpu', 'memory', 'ready', 'status',];
  private readonly apiService = inject(ApiService);
  private readonly router = inject(Router);
  dataSource: Signal<Pod[]>;

  constructor() {
    this.dataSource = toSignal(this.#refresh.pipe(
      switchMap(() => this.apiService.getPods())
    ), {initialValue: []});
  }

  refresh(): void {
    this.#refresh.next(undefined);
  }

  showLogs(pod: Pod) {
    void this.router.navigate(['pods', 'logs', pod.namespace, pod.name])
  }
}
