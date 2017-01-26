/// <reference path="../References.d.ts"/>
import * as SuperAgent from 'superagent';
import Dispatcher from '../dispatcher/Dispatcher';
import * as Alert from '../Alert';
import * as BuildTypes from '../types/BuildTypes';

export function sync(): Promise<string> {
	Dispatcher.dispatch({
		type: BuildTypes.LOADING,
	});

	return new Promise<string>((resolve, reject): void => {
		SuperAgent
			.get('/build')
			.set('Accept', 'application/json')
			.end((err: any, res: SuperAgent.Response): void => {
				Dispatcher.dispatch({
					type: BuildTypes.LOADED,
				});

				if (err) {
					Alert.error(err);
					reject(err);
					return;
				}

				Dispatcher.dispatch({
					type: BuildTypes.SYNC,
					data: {
						builds: res.body,
					},
				});

				resolve();
			});
	});
}

export function remove(id: string): void {
	Dispatcher.dispatch({
		type: BuildTypes.REMOVE,
		data: {
			id: id,
		},
	});

	Alert.info('Build successfully removed');
}
