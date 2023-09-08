using test from '../db/data-model';

service TestService {
  entity Temperatures @readonly as projection on test.Temperatures;
}
