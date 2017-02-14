/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import Loader from '../Loader';
import * as BuildInfoTypes from '../types/BuildInfoTypes';

export function open(id: string): Promise<void> {
	let loader = new Loader().loading();

	Dispatcher.dispatch({
		type: BuildInfoTypes.OPEN,
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
					type: BuildInfoTypes.UPDATE,
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
		type: BuildInfoTypes.CLOSE,
	});
}