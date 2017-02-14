/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as BuildLogTypes from '../types/BuildLogTypes';

export function open(id: string): Promise<void> {
	let loader = new Loader().loading();

	Dispatcher.dispatch({
		type: BuildLogTypes.OPEN,
		data: {
			id: id,
		},
	});

	return new Promise<void>((resolve, reject): void => {
		SuperAgent
			.get('/build/' + id + '/log')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				loader.done();

				if (err) {
					Alert.error('Failed to get build log');
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: BuildLogTypes.UPDATE,
					data: {
						output: res.body,
					},
				});

				resolve();
			});
	});
}

export function close(): void {
	Dispatcher.dispatch({
		type: BuildLogTypes.CLOSE,
	});
}
