/// <reference path="References.d.ts"/>
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as injectTapEventPlugin from 'react-tap-event-plugin';
import Main from './components/Main';

injectTapEventPlugin();

ReactDOM.render(
	<Main/>,
	document.getElementById('app')
);
