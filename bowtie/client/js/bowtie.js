
/////////////////////////////////////////////////////////
//
// Bowtie bindings for JavaScript
// 
// Github  : http://github.com/wallarelvo/Bowtie
// Contact : aw204@st-andrews.ac.uk
// Website : www.bowtie.mobi
//
/////////////////////////////////////////////////////////

/* * * * * * * * * * * * * * * * * * * * * * * * * * * */

// namespaces
bowtie = {};
bowtie.util = {};
bowtie.prefixes = {};
bowtie.constants = {};

/////////////////////////////////////////////////////////
//
// prefixes
//
// Url prefixes used in bowtie
//
/////////////////////////////////////////////////////////

bowtie.prefixes.MEDIA = "media";
bowtie.prefixes.SENSORS = "sensors";
bowtie.prefixes.NODES = "nodes";

/////////////////////////////////////////////////////////
//
// constants
//
// Constants used in bowtie
//
/////////////////////////////////////////////////////////

bowtie.constants.AUDIO = "audio";
bowtie.constants.VIDEO = "video";
bowtie.constants.BOWTIE_URL = "www.bowtie.mobi";

/////////////////////////////////////////////////////////
//
// util.Callback
//
// Class used to organize callback functions 
//
/////////////////////////////////////////////////////////

bowtie.util.Callback = function (url, prefix, groupId, nodeId, sensor) {
	this.path = url + "/" + prefix + "/";
	if (groupId != undefined) {
		this.path += groupId + "/";
		if (nodeId != undefined) {
			this.path += nodeId + "/";
			if (sensor != undefined) {
				this.path += sensor + "/";
			}
		}
	}
}

bowtie.util.Callback.prototype = {
	setCallback : function (callback) {
		$.getJSON(this.path, callback);
	}
}

/////////////////////////////////////////////////////////
//
// BowtieClient 
//
// Class used for making requests to the bowtie server 
// of your choosing
//
/////////////////////////////////////////////////////////

bowtie.BowtieClient = function (url) {
	this.url = url ? url : bowtie.constants.BOWTIE_URL;

	if (this.url.slice(0, 4) != "http") {
		this.url = "http://" + this.url;
	}
}

bowtie.BowtieClient.prototype = {
	getSensor : function (groupId, nodeId, sensor) {
		return new bowtie.util.Callback(
			this.url,
			bowtie.prefixes.SENSORS,
			groupId,
			nodeId,
			sensor
		);
	}, 

	getNode : function (groupId, nodeId) {
		return new bowtie.util.Callback(
			this.url,
			bowtie.prefixes.SENSORS,
			groupId, 
			nodeId
		);
	},

	getGroup : function (groupId) {
		return new bowtie.util.Callback(
			this.url,
			bowtie.prefixes.SENSORS,
			groupId
		);
	},

	getNodeList : function (groupId) {
		return new bowtie.util.Callback(
			this.url,
			bowtie.prefixes.NODES,
			groupId
		);
	},

	deleteSensor : function (groupId, nodeId, sensor) {
		return $.ajax(
			{
				type : "DELETE",

				url : (
					this.url + "/" + 
					bowtie.prefixes.SENSORS + "/" + 
					groupId + "/" + 
					nodeId + "/" + 
					sensor
				)
			}
		);
	}, 

	deleteNode : function (groupId, nodeId) {
		return $.ajax(
			{
				type : "DELETE",

				url : (
					this.url + "/" +
					bowtie.prefixes.SENSORS + "/" + 
					groupId + "/" + 
					nodeId
				)
			}
		);
	},

	deleteGroup : function (groupId, nodeId) {
		return $.ajax(
			{
				type : "DELETE",

				url : (
					this.url + "/" +
					bowtie.prefixes.SENSORS + "/" + 
					groupId
				)
			}
		);
	}
}




