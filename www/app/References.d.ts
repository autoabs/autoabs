declare namespace JSX {
	interface IntrinsicElements {
		'iron-a11y-announcer': any,
		'iron-a11y-keys-behavior': any,
		'iron-autogrow-textarea': any,
		'iron-behaviors': any,
		'iron-checked-element-behavior': any,
		'iron-collapse': any,
		'iron-dropdown': any,
		'iron-fit-behavior': any,
		'iron-flex-layout': any,
		'iron-form-element-behavior': any,
		'iron-icon': any,
		'iron-icons': any,
		'iron-iconset-svg': any,
		'iron-input': any,
		'iron-media-query': any,
		'iron-menu-behavior': any,
		'iron-meta': any,
		'iron-overlay-behavior': any,
		'iron-range-behavior': any,
		'iron-resizable-behavior': any,
		'iron-selector': any,
		'iron-validatable-behavior': any,
		'neon-animation': any,
		'paper-badge': any,
		'paper-behaviors': any,
		'paper-button': any,
		'paper-card': any,
		'paper-checkbox': any,
		'paper-dialog': any,
		'paper-dialog-behavior': any,
		'paper-dialog-scrollable': any,
		'paper-drawer-panel': any,
		'paper-dropdown-menu': any,
		'paper-elements': any,
		'paper-fab': any,
		'paper-header-panel': any,
		'paper-icon-button': any,
		'paper-input': any,
		'paper-item': any,
		'paper-material': any,
		'paper-menu': any,
		'paper-menu-button': any,
		'paper-progress': any,
		'paper-radio-button': any,
		'paper-radio-group': any,
		'paper-ripple': any,
		'paper-scroll-header-panel': any,
		'paper-slider': any,
		'paper-spinner': any,
		'paper-styles': any,
		'paper-tabs': any,
		'paper-toast': any,
		'paper-toggle-button': any,
		'paper-toolbar': any,
		'paper-tooltip': any,
		'google-analytics': any,
		'google-analytics-chart': any,
		'google-analytics-dashboard': any,
		'google-analytics-date-selector': any,
		'google-analytics-loader': any,
		'google-analytics-query': any,
		'google-analytics-view-selector': any,
		'google-calendar': any,
		'google-castable-video': any,
		'google-chart': any,
		'google-feeds': any,
		'google-hangout-button': any,
		'google-map': any,
		'google-map-directions': any,
		'google-map-elements': any,
		'google-map-marker': any,
		'google-map-point': any,
		'google-map-poly': any,
		'google-map-search': any,
		'google-sheets': any,
		'google-signin': any,
		'google-streetview-pano': any,
		'google-youtube': any,
		'google-youtube-upload': any,
		're-captcha': any,
	}
}

declare module 'chartjs' {
	class Chart {
		constructor(ctx: CanvasRenderingContext2D, options?: any);

		static Line(ctx: CanvasRenderingContext2D, options?: any): Chart;

		static Bar(ctx: CanvasRenderingContext2D, options?: any): Chart;

		static Bubble(ctx: CanvasRenderingContext2D, options?: any): Chart;

		static Doughnut(ctx: CanvasRenderingContext2D, options?: any): Chart;

		static PolarArea(ctx: CanvasRenderingContext2D, options?: any): Chart;

		static Radar(ctx: CanvasRenderingContext2D, options?: any): Chart;
	}

	export interface LinearChartData {
		labels: string[];
		datasets: ChartDataSet[];
		xAxisID?: string;
		yAxisID?: string;
		fill?: boolean;
		lineTension?: number;
		backgroundColor?: string;
		borderWidth?: number;
		borderColor?: string;
		borderCapStyle?: string;
		borderDash?: number[];
		borderDashOffset?: number;
		borderJoinStyle?: string;
		pointBorderColor?: string[];
		pointBorderWidth?: number[];
		pointRadius?: number[];
		pointHoverRadius?: number[];
		pointHitRadius?: number[];
		pointHoverBackgroundColor?: string[];
		pointHoverBorderColor?: string[];
		pointHoverBorderWidth?: number[];
		pointStyle?: any;
	}

	export interface CircularChartData {
		value: number;
		color?: string;
		highlight?: string;
		label?: string;
	}

	export interface ChartDataSet {
		label?: string;
		fillColor?: string;
		strokeColor?: string;
		borderColor?: string | string[];
		backgroundColor?: string | string[];
		borderWidth?: number;
		pointColor?: string;
		pointStrokeColor?: string;
		pointHighlightFill?: string;
		pointHighlightStroke?: string;
		highlightFill?: string;
		highlightStroke?: string;
		data: number[];
	}

	export interface ChartSettings {
		animation?: boolean;
		animationSteps?: number;
		animationEasing?: string;
		showScale?: boolean;
		scaleOverride?: boolean;
		scaleSteps?: number;
		scaleStepWidth?: number;
		scaleStartValue?: number;
		scaleLineColor?: string;
		scaleLineWidth?: number;
		scaleShowLabels?: boolean;
		scaleLabel?: string;
		scaleIntegersOnly?: boolean;
		scaleBeginAtZero?: boolean;
		scaleFontFamily?: string;
		scaleFontSize?: number;
		scaleFontStyle?: string;
		scaleFontColor?: string;
		responsive?: boolean;
		maintainAspectRatio?: boolean;
		showTooltips?: boolean;
		tooltipEvents?: string[];
		tooltipFillColor?: string;
		tooltipFontFamily?: string;
		tooltipFontSize?: number;
		tooltipFontStyle?: string;
		tooltipFontColor?: string;
		tooltipTitleFontFamily?: string;
		tooltipTitleFontSize?: number;
		tooltipTitleFontStyle?: string;
		tooltipTitleFontColor?: string;
		tooltipYPadding?: number;
		tooltipXPadding?: number;
		tooltipCaretSize?: number;
		tooltipCornerRadius?: number;
		tooltipXOffset?: number;
		tooltipTemplate?: string;
		multiTooltipTemplate?: string;
		onAnimationProgress?: () => any;
		onAnimationComplete?: () => any;
	}

	export interface ChartOptions extends ChartSettings {
		scaleShowGridLines?: boolean;
		scaleGridLineColor?: string;
		scaleGridLineWidth?: number;
		scaleShowHorizontalLines?: boolean;
		scaleShowVerticalLines?: boolean;
		legendTemplate?: string;
	}

	export interface PointsAtEvent {
		value: number;
		label: string;
		datasetLabel: string;
		strokeColor: string;
		fillColor: string;
		highlightFill: string;
		highlightStroke: string;
		x: number;
		y: number;
	}

	export interface ChartInstance {
		clear: () => void;
		stop: () => void;
		resize: () => void;
		destroy: () => void;
		toBase64Image: () => string;
		generateLegend: () => string;
	}

	export interface LinearInstance extends ChartInstance {
		getPointsAtEvent: (event: Event) => PointsAtEvent[];
		update: () => void;
		addData: (valuesArray: number[], label: string) => void;
		removeData: (index?: number) => void;
	}

	export interface CircularInstance extends ChartInstance {
		getSegmentsAtEvent: (event: Event) => {}[];
		update: () => void;
		addData: (valuesArray: CircularChartData, index?: number) => void;
		removeData: (index: number) => void;
		segments: Array<CircularChartData>;
	}

	export interface LineChartOptions extends ChartOptions {
		bezierCurve?: boolean;
		bezierCurveTension?: number;
		pointDot?: boolean;
		pointDotRadius?: number;
		pointDotStrokeWidth?: number;
		pointHitDetectionRadius?: number;
		datasetStroke?: boolean;
		datasetStrokeWidth?: number;
		datasetFill?: boolean;
	}

	export interface BarChartOptions extends ChartOptions {
		scaleBeginAtZero?: boolean;
		barShowStroke?: boolean;
		barStrokeWidth?: number;
		barValueSpacing?: number;
		barDatasetSpacing?: number;
	}

	export interface RadarChartOptions extends ChartSettings {
		scaleShowLine?: boolean;
		angleShowLineOut?: boolean;
		scaleShowLabels?: boolean;
		scaleBeginAtZero?: boolean;
		angleLineColor?: string;
		angleLineWidth?: number;
		pointLabelFontFamily?: string;
		pointLabelFontStyle?: string;
		pointLabelFontSize?: number;
		pointLabelFontColor?: string;
		pointDot?: boolean;
		pointDotRadius?: number;
		pointDotStrokeWidth?: number;
		pointHitDetectionRadius?: number;
		datasetStroke?: boolean;
		datasetStrokeWidth?: number;
		datasetFill?: boolean;
		legendTemplate?: string;
	}

	export interface PolarAreaChartOptions extends ChartSettings {
		scaleShowLabelBackdrop?: boolean;
		scaleBackdropColor?: string;
		scaleBeginAtZero?: boolean;
		scaleBackdropPaddingY?: number;
		scaleBackdropPaddingX?: number;
		scaleShowLine?: boolean;
		segmentShowStroke?: boolean;
		segmentStrokeColor?: string;
		segmentStrokeWidth?: number;
		animationSteps?: number;
		animationEasing?: string;
		animateRotate?: boolean;
		animateScale?: boolean;
		legendTemplate?: string;
	}

	export interface PieChartOptions extends ChartSettings {
		segmentShowStroke?: boolean;
		segmentStrokeColor?: string;
		segmentStrokeWidth?: number;
		percentageInnerCutout?: number;
		animationSteps?: number;
		animationEasing?: string;
		animateRotate?: boolean;
		animateScale?: boolean;
		legendTemplate?: string;
	}
}
