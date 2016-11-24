/// <reference path="../References.d.ts"/>
import * as React from 'react';
import * as BuildTypes from '../types/BuildTypes';
import * as MiscUtils from '../utils/MiscUtils';

interface Props {
	build: BuildTypes.Build;
}

const css = {
	box: {
		display: 'table-row',
	} as React.CSSProperties,
	field: {
		display: 'table-cell',
	} as React.CSSProperties,
};

export default class Build extends React.Component<Props, null> {
	render(): JSX.Element {
		let build = this.props.build;

		let start = '-';
		if (build.start !== "0001-01-01T00:00:00Z") {
			console.log(build.start);
			start = MiscUtils.formatDate(new Date(build.start));
		}
		let stop = '-';
		if (build.stop !== "0001-01-01T00:00:00Z") {
			stop = MiscUtils.formatDate(new Date(build.stop));
		}

		return <div style={css.box}>
			<div style={css.field}>{build.id}</div>
			<div style={css.field}>{build.name}</div>
			<div style={css.field}>{build.state}</div>
			<div style={css.field}>{build.version}</div>
			<div style={css.field}>{build.release}</div>
			<div style={css.field}>{build.repo}</div>
			<div style={css.field}>{build.arch}</div>
			<div style={css.field}>{start}</div>
			<div style={css.field}>{stop}</div>
		</div>;
	}
}
