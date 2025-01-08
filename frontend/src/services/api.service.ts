import {Injectable} from '@angular/core';
import {catchError, firstValueFrom, map, Observable, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {ContainerDeployment} from '../models/container-deployment.model';
import {Pod} from '../models/pod.model';

@Injectable({providedIn: 'root'})
export class ApiService {

  baseUrl = 'http://localhost:5195/';

  constructor(private readonly http: HttpClient) {
  }

  getDeployments(): Observable<ContainerDeployment[]> {
    return this.http.get<ContainerDeployment[]>(this.baseUrl + 'containerDeployments');
  }

  getNamespaces(): Observable<string[]> {
    return this.http.get<string[]>(this.baseUrl + 'namespaces');
  }

  getPods(): Observable<Pod[]> {
    return this.http.get<Pod[]>(this.baseUrl + 'pods');
  }

  createDeployment(request: ContainerDeployment): Promise<boolean> {
    return firstValueFrom(this.http.post(this.baseUrl + 'createDeployment', request).pipe(
      map(() => true),
      catchError(() => of(false))
    ));
  }

  updateDeployment(request: ContainerDeployment): Promise<boolean> {
    return firstValueFrom(this.http.post(this.baseUrl + 'updateDeployment', request).pipe(
      map(() => true),
      catchError(() => of(false))
    ));
  }
}
