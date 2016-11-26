/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import * as BuildActions from '../actions/BuildActions';

let loaded = false;

export function load(): Promise<string> {
	if (!loaded) {
		BuildActions.loading();
	}
	loaded = true;

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.get('/builds')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				if (err) {
					reject(err);
					return;
				}

				BuildActions.load(res.body);
				resolve();
			})
	});
}
