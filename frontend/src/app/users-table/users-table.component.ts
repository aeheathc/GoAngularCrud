import { Component } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { SettingsService } from '../settings/settings.service';

@Component({
  selector: 'app-users-table',
  templateUrl: './users-table.component.html',
  styleUrls: ['./users-table.component.css']
})
export class UsersTableComponent {
  users: User[] = [];
  backendUrl: string;

  constructor(private http: HttpClient, private settings: SettingsService) {
    this.backendUrl = this.settings.settings.backendUrl ?? 'http://localhost:8888';
  }

  ngOnInit() {
    this.getUsers();
  }

  postUser(userName: string, firstName: string, lastName: string, email: string, status: string, department: string, newUserForm: any) {
    console.log("post starting");
    const headers = new HttpHeaders().set('Content-Type', 'application/json; charset=utf-8');
    this.http.post(
      this.backendUrl+'/users',
      new User(userName, firstName, lastName, email, status, department),
      {headers: headers}
    ).subscribe(_ => {
      this.getUsers();
      newUserForm.reset();
    });
  }

  getUsers() {
    this.http.get(this.backendUrl+'/users').subscribe(data => {
      this.users = [];
      if(data != null) {
        for(let user of data as User[]) {
          this.users.push(user);
        }
      }
    })
  }

  deleteUser(id: number | undefined) {
    this.http.delete(this.backendUrl+'/users/'+id).subscribe(_ => {
      this.getUsers();
    });
  }

  saveUser(id: number | undefined) {
    const headers = new HttpHeaders().set('Content-Type', 'application/json; charset=utf-8');
    this.http.put(
      this.backendUrl+'/users/'+id,
      new User(
        (<HTMLInputElement>document.getElementById('user'+id+'username')).value,
        (<HTMLInputElement>document.getElementById('user'+id+'firstname')).value,
        (<HTMLInputElement>document.getElementById('user'+id+'lastname')).value,
        (<HTMLInputElement>document.getElementById('user'+id+'email')).value,
        (<HTMLInputElement>document.getElementById('user'+id+'status')).value,
        (<HTMLInputElement>document.getElementById('user'+id+'department')).value
      ),
      {headers: headers}
    ).subscribe(_ => {
      this.getUsers();
    });
  }
}

class User
{
  user_id?: number
  user_name: string
  first_name: string
  last_name: string
  email: string
  user_status: string
  department: string

  constructor(user_name: string, first_name: string, last_name: string, email: string, user_status: string, department: string, user_id?: number) {
    this.user_id = user_id;
    this.user_name = user_name;
    this.first_name = first_name;
    this.last_name = last_name;
    this.email = email;
    this.user_status = user_status;
    this.department = department;
  }
}