import {Component, HostBinding, inject, Signal} from '@angular/core';
import {
  MAT_DIALOG_DATA,
  MatDialogActions,
  MatDialogClose,
  MatDialogContent,
  MatDialogRef,
  MatDialogTitle
} from '@angular/material/dialog';
import {MatButton, MatIconButton} from '@angular/material/button';
import {FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators} from '@angular/forms';
import {ApiService} from '../../services/api.service';
import {toSignal} from '@angular/core/rxjs-interop';
import {MatFormField, MatHint, MatLabel} from '@angular/material/form-field';
import {MatOption, MatSelect} from '@angular/material/select';
import {MatInput} from '@angular/material/input';
import {ContainerDeployment} from '../../models/container-deployment.model';
import {MatIcon} from '@angular/material/icon';

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
    MatIconButton,
    MatIcon,
  ],
  templateUrl: './create-deployment-dialog.component.html',
})
export class CreateDeploymentDialogComponent {
  readonly namespaceOptions: Signal<string[]> = toSignal(inject(ApiService).getNamespaces(), {initialValue: []});
  private readonly apiService = inject(ApiService);
  private readonly dialogRef = inject(MatDialogRef);
  readonly data?: ContainerDeployment = inject(MAT_DIALOG_DATA);


  private readonly formBuilder = inject(FormBuilder);
  form: FormGroup = this.formBuilder.group({
    namespace: ['', Validators.required],
    name: ['', Validators.required],
    image: ['', Validators.required],
    cpu: ['', Validators.required],
    memory: ['', Validators.required],
    port: ['', Validators.required],
    replicas: [1, Validators.required],
    environmentVariables: this.formBuilder.array([])
  });

  constructor() {
    if (!this.data) {
      return;
    }

    const {environmentVariables, ...rest} = this.data;
    if (this.data) {
      this.form.patchValue(rest);
      this.form.controls['namespace'].disable();
      this.form.controls['name'].disable();
      this.form.controls['image'].disable();
      if (environmentVariables) {
        const envs = Object.entries(environmentVariables);
        envs.forEach(([key, value]) => this.addEnvVariable(key, value));
      }
    }
  }

  get envVariables(): FormArray {
    return this.form.get('environmentVariables') as FormArray;
  }

  addEnvVariable(key: string = '', value: string = ''): void {
    const envGroup = this.formBuilder.group({
      key: [key, Validators.required],
      value: [value, Validators.required],
    });
    this.envVariables.push(envGroup);
  }

  removeEnvVariable(index: number): void {
    this.envVariables.removeAt(index);
  }


  async save(): Promise<void> {
    if (!this.form.valid) {
      return;
    }

    const {environmentVariables, ...rest} = this.form.getRawValue();
    const mappedEnv = environmentVariables.reduce((acc: any, {key, value}: any) => ({...acc, [key]: value}), {});
    if (!this.data) {
      if (await this.apiService.createDeployment({...rest, environmentVariables: mappedEnv})) {
        this.dialogRef.close(true)
      }
      return;
    }

    if (await this.apiService.updateDeployment({...rest, environmentVariables: mappedEnv})) {
      this.dialogRef.close(true)
    }
  }
}
