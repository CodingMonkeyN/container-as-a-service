import {Component, HostBinding, inject, Signal} from '@angular/core';
import {
  MatCell, MatCellDef,
  MatColumnDef,
  MatHeaderCell,
  MatHeaderCellDef,
  MatHeaderRow, MatHeaderRowDef,
  MatRow, MatRowDef,
  MatTable
} from '@angular/material/table';
import {MatButton, MatIconButton} from '@angular/material/button';
import {ApiService} from '../../services/api.service';
import {toSignal} from '@angular/core/rxjs-interop';
import {MatCard, MatCardContent, MatCardHeader, MatCardTitle} from '@angular/material/card';
import {MatDialog} from '@angular/material/dialog';
import {CreateDeploymentDialogComponent} from '../create-deployment-dialog/create-deployment-dialog.component';
import {BehaviorSubject, switchMap, tap} from 'rxjs';
import {ContainerDeployment} from '../../models/container-deployment.model';
import {MatIcon} from '@angular/material/icon';

@Component({
  selector: 'app-deployments',
  imports: [
    MatTable,
    MatColumnDef,
    MatHeaderCell,
    MatCell,
    MatHeaderRow,
    MatRow,
    MatHeaderCellDef,
    MatCellDef,
    MatButton,
    MatHeaderRowDef,
    MatRowDef,
    MatCard,
    MatCardTitle,
    MatCardHeader,
    MatCardContent,
    MatIconButton,
    MatIcon
  ],
  templateUrl: './deployments.component.html',
})
export class DeploymentsComponent {
  @HostBinding() class = 'flex-1 flex flex-col h-full w-full';

  readonly #refresh = new BehaviorSubject<void>(undefined);

  displayedColumns: string[] = ['namespace', 'name', 'image', 'port', 'cpu', 'memory', 'replicas', 'actions'];
  dataSource: Signal<ContainerDeployment[]>;
  private readonly dialog = inject(MatDialog);
  private readonly api = inject(ApiService);

  constructor() {
    this.dataSource = toSignal(this.#refresh.pipe(
      tap(() => console.log('Refreshing deployments')),
      switchMap(() => this.api.getDeployments())
    ), {initialValue: []});
  }

  addDeployment(): void {
    this.dialog.open(CreateDeploymentDialogComponent);
  }

  editDeployment(deployment: ContainerDeployment) {
    const ref = this.dialog.open(CreateDeploymentDialogComponent, {data: deployment,});
    ref.afterClosed().subscribe(result => {
      if (result === true) {
        this.#refresh.next();
      }
    })
  }

  refresh() {
    this.#refresh.next();
  }
}
