import {Component, inject, Signal} from '@angular/core';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from '@angular/material/dialog';
import {MatButton} from '@angular/material/button';
import {FormBuilder, FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {ApiService} from '../../services/api.service';
import {toSignal} from '@angular/core/rxjs-interop';
import {MatFormField, MatHint, MatLabel} from '@angular/material/form-field';
import {MatOption, MatSelect} from '@angular/material/select';
import {MatInput} from '@angular/material/input';
import {ContainerDeployment} from '../../models/container-deployment.model';

@Component({
  selector: 'app-create-deployment-dialog',
  imports: [
    MatDialogTitle,
    ReactiveFormsModule,
    MatDialogContent,
    MatDialogActions,
    MatButton,
    MatDialogClose,
    ReactiveFormsModule,
    MatFormField,
    MatSelect,
    MatLabel,
    MatOption,
    MatInput,
    MatHint,
  ],
  templateUrl: './create-deployment-dialog.component.html',
})
export class CreateDeploymentDialogComponent {

  readonly namespaceOptions: Signal<string[]> = toSignal(inject(ApiService).getNamespaces(), {initialValue: []});
  private readonly apiService = inject(ApiService);
  private readonly dialogRef = inject(MatDialogRef);
  readonly data?: ContainerDeployment = inject(MAT_DIALOG_DATA);


  form: FormGroup = inject(FormBuilder).group({
    namespace: ['', Validators.required],
    name: ['', Validators.required],
    image: ['', Validators.required],
    cpu: ['', Validators.required],
    memory: ['', Validators.required],
    port: ['', Validators.required],
    replicas: [1, Validators.required],
  });

  constructor() {
    if (this.data) {
      this.form.patchValue(this.data);
      this.form.controls['namespace'].disable();
      this.form.controls['name'].disable();
      this.form.controls['image'].disable();
    }
  }

  async save(): Promise<void> {
    if (!this.form.valid) {
      return;
    }

    if (!this.data) {
      if (await this.apiService.createDeployment(this.form.getRawValue())) {
        this.dialogRef.close(true)
      }
      return;
    }

    if (await this.apiService.updateDeployment(this.form.getRawValue())) {
      this.dialogRef.close(true)
    }
  }
}
