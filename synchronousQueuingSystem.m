close all
clear, clc

figure()
ideal = load('res.txt');
xi = ideal(:,3); 
yi = ideal(:,1);
xlabel('Lambda')
ylabel('N')
ylim([0 5])
plot(xi, yi); 
grid on;
hold on;

ideal = load('resTheor.txt');
xi_t = ideal(:,3); 
yi_t = ideal(:,1);
plot(xi_t, yi_t); 
legend('Practic','Theor')
title('N asynchronous');



figure()
ideal = load('res.txt');
xi = ideal(:,3); 
yi = ideal(:,2);
plot(xi, yi); 
ylim([0 6])
xlabel('Lambda')
ylabel('D')
grid on;
hold on;
ideal = load('resTheor.txt');
xi_t = ideal(:,3); 
yi_t = ideal(:,2);
plot(xi_t, yi_t); 
legend('Practic','Theor')
title('D asynchronous');






