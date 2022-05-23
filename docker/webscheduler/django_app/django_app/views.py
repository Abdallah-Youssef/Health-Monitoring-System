from django.shortcuts import render,redirect
from .forms import QueryForm
from django.http import HttpResponse
import datetime
from urllib.parse import urlencode
from django.urls import reverse

def query_result_view(request,start_date,end_date):
    start_date_splitted = start_date.split()
    start_date = start_date_splitted[0]+" "+start_date_splitted[1]
    end_date = start_date_splitted[2]+" "+end_date

    return HttpResponse("<h1> service 1</h1>")


def query_form_view(request):
    # submit button is pushed
    if request.method == 'POST':
        # extract start time and end time
        start_date = request.POST['start_date']
        end_date = request.POST['end_date']
        # get the form to validate the inputs
        form = QueryForm(request.POST)
        if form.is_valid():
            # convert into datetime to make sure that end time > start time
            start_date_obj = datetime.datetime.strptime(start_date, '%d/%m/%Y %H:%M')
            end_date_obj = datetime.datetime.strptime(end_date, '%d/%m/%Y %H:%M')
            if end_date_obj < start_date_obj:
                # end time is before start time so it will preview an error msg
                context = {'form': QueryForm(), 'error_msg': "end date cannot be before start date!"}
                return render(request, "queryForm.html", context)
            else:
                #inputs are valid
                # base_url = reverse('query')
                # print(base_url)
                # query_string=urlencode({'start_date':str(start_date),'end_date':str(end_date)})
                # print(query_string)
                # url = '{}?{}'.format(base_url, query_string)
                return redirect("query",start_date=str(start_date),end_date=str(end_date))
        else:
            context = {'form': QueryForm(), 'error_msg': "follow the input format: d/m/Y H:M !"}
            return render(request, "queryForm.html", context)
    #calling the form page
    elif request.method == 'GET':
        context = {'form': QueryForm()}
        return render(request, "queryForm.html", context)


