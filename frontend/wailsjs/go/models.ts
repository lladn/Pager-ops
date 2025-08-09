export namespace database {
	
	export class CachedIncident {
	    id: string;
	    service_id: string;
	    status: string;
	    title: string;
	    description: string;
	    urgency: string;
	    incident_number: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	    pinned: boolean;
	    assigned_user_ids: string[];
	    escalation_level: number;
	    html_url: string;
	
	    static createFrom(source: any = {}) {
	        return new CachedIncident(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.service_id = source["service_id"];
	        this.status = source["status"];
	        this.title = source["title"];
	        this.description = source["description"];
	        this.urgency = source["urgency"];
	        this.incident_number = source["incident_number"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	        this.pinned = source["pinned"];
	        this.assigned_user_ids = source["assigned_user_ids"];
	        this.escalation_level = source["escalation_level"];
	        this.html_url = source["html_url"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class DraftNote {
	    incident_id: string;
	    note_text: string;
	    why_triggered: string;
	    impact: string;
	    actions: string;
	    links: string;
	    // Go type: time
	    last_updated: any;
	
	    static createFrom(source: any = {}) {
	        return new DraftNote(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.incident_id = source["incident_id"];
	        this.note_text = source["note_text"];
	        this.why_triggered = source["why_triggered"];
	        this.impact = source["impact"];
	        this.actions = source["actions"];
	        this.links = source["links"];
	        this.last_updated = this.convertValues(source["last_updated"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Service {
	    id: string;
	    alias: string;
	    enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.alias = source["alias"];
	        this.enabled = source["enabled"];
	    }
	}
	export class Template {
	    id: number;
	    title: string;
	    body_text: string;
	
	    static createFrom(source: any = {}) {
	        return new Template(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.body_text = source["body_text"];
	    }
	}
	export class User {
	    id: string;
	    email: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.email = source["email"];
	        this.name = source["name"];
	    }
	}

}

export namespace main {
	
	export class AppSettings {
	    api_key: string;
	    theme: string;
	    assigned_only: boolean;
	    redirect_enabled: boolean;
	    sound_path: string;
	    sound_enabled: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.api_key = source["api_key"];
	        this.theme = source["theme"];
	        this.assigned_only = source["assigned_only"];
	        this.redirect_enabled = source["redirect_enabled"];
	        this.sound_path = source["sound_path"];
	        this.sound_enabled = source["sound_enabled"];
	    }
	}

}

export namespace services {
	
	export class Alert {
	    id: string;
	    summary: string;
	    status: string;
	    created_at: string;
	    details: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new Alert(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.summary = source["summary"];
	        this.status = source["status"];
	        this.created_at = source["created_at"];
	        this.details = source["details"];
	    }
	}
	export class Service {
	    id: string;
	    name: string;
	    description: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Service(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.status = source["status"];
	    }
	}

}

