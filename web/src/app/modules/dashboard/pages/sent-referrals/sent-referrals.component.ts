import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params, RouterLink } from '@angular/router';
import { catchError, of, switchMap } from 'rxjs';
import { EmailingService, ReferralMailbox } from '../../services/emailing.service';
import { AsyncPipe, DatePipe, TitleCasePipe } from '@angular/common';

@Component({
  selector: 'app-sent-referrals',
  imports: [RouterLink, TitleCasePipe, DatePipe],
  templateUrl: './sent-referrals.component.html',
  styleUrl: './sent-referrals.component.css'
})
export class SentReferralsComponent implements OnInit {
  apiErrorMessage: string | null = null;
  referralEmails: ReferralMailbox[] = [];
  constructor(private readonly route: ActivatedRoute, private readonly emailService: EmailingService) { }

  ngOnInit(): void {
    this.route.params.pipe(
      switchMap((param: Params) =>
        this.emailService.getReferralEmailByUuid$(param["uuid"]
        )
      ), catchError(err => {
        this.apiErrorMessage = JSON.stringify(err.error) || 'Something went wrong!'
        return of(null)
      })).subscribe(data => {
        if (data === null) return;
        this.referralEmails = data as ReferralMailbox[];
        this.apiErrorMessage = null;
        console.log(data);
      })
  }
}
