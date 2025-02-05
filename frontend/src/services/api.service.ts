import {Injectable} from '@angular/core';
import {catchError, firstValueFrom, map, Observable, of} from 'rxjs';
import {HttpClient} from '@angular/common/http';
import {ContainerDeployment} from '../models/container-deployment.model';
import {Pod} from '../models/pod.model';

@Injectable({providedIn: 'root'})
export class ApiService {

  baseUrl = 'https://localhost:5001/api/';

  constructor(private readonly http: HttpClient) {
  }

  getDeployments(): Observable<ContainerDeployment[]> {
    return this.http.get<ContainerDeployment[]>(this.baseUrl + 'containerDeployment');
  }

  getNamespaces(): Observable<string[]> {
    return this.http.get<string[]>(this.baseUrl + 'namespace');
  }

  getPods(): Observable<Pod[]> {
    return this.http.get<Pod[]>(this.baseUrl + 'pods');
  }

  getLogs(namespace: string, containerName: string): Observable<string> {
    return new Observable<string>(observer => {
      const xhr = new XMLHttpRequest();
      xhr.open('GET', this.baseUrl + 'pods/logs/' + namespace + '/' + containerName, true);
      xhr.responseType = 'arraybuffer';

      xhr.onload = () => {
        if (xhr.status === 200) {
          const decodedLogs = new TextDecoder('utf-8').decode(xhr.response);
          console.log(decodedLogs)
          observer.next(decodedLogs);
        }
      };

      xhr.onerror = () => observer.error('Fehler beim Laden der Logs.');
      xhr.send();

      return () => xhr.abort();
    });
  }

  createDeployment(request: ContainerDeployment): Promise<boolean> {
    return firstValueFrom(this.http.post(this.baseUrl + 'containerDeployment/create', request).pipe(
      map(() => true),
      catchError(() => of(false))
    ));
  }

  updateDeployment(request: ContainerDeployment): Promise<boolean> {
    return firstValueFrom(this.http.post(this.baseUrl + 'containerDeployment/update', request).pipe(
      map(() => true),
      catchError(() => of(false))
    ));
  }
}
